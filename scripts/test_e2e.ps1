param(
    [string]$BaseUrl = 'http://localhost:18080',
    [string]$DatabaseContainer = 'wallet-qa-postgres',
    [string]$DatabaseName = 'qa_wallet'
)
$ErrorActionPreference = 'Stop'
$script:checks = 0

function Request($Name, $Method, $Path, $Expected, $Body = $null, $Token = '') {
    $params = @{ Uri = "$BaseUrl$Path"; Method = $Method; SkipHttpErrorCheck = $true; TimeoutSec = 15 }
    if ($null -ne $Body) { $params.Body = $Body | ConvertTo-Json; $params.ContentType = 'application/json' }
    if ($Token) { $params.Headers = @{ Authorization = "Bearer $Token" } }
    $response = Invoke-WebRequest @params
    if ([int]$response.StatusCode -ne $Expected) { throw "$Name expected HTTP $Expected, got $($response.StatusCode): $($response.Content)" }
    $script:checks++
    Write-Host "PASS: $Name (HTTP $Expected)"
    if ($response.Headers.'Content-Type' -match 'application/json') { return $response.Content | ConvertFrom-Json }
}
function AssertBalance($Token, $Expected) {
    $response = Request 'balance' GET '/api/v1/wallet/balance' 200 $null $Token
    if ([decimal]$response.data.balance -ne $Expected) { throw "Balance expected $Expected, got $($response.data.balance)" }
    $script:checks++
}

$null = Request 'health' GET '/health' 200
$null = Request 'Swagger UI' GET '/swagger/index.html' 200
$null = Request 'Swagger specification' GET '/swagger/doc.json' 200
$null = Request 'anonymous access rejected' GET '/api/v1/wallet/balance' 401
$null = Request 'invalid token rejected' GET '/api/v1/auth/me' 401 $null 'invalid'
$suffix = [guid]::NewGuid().ToString('N')
$email1 = "qa-$suffix@example.com"
$email2 = "qa-recipient-$suffix@example.com"
$user1 = Request 'register sender' POST '/api/v1/auth/register' 201 @{name='QA Sender';cpf='12345678909';email=$email1;phone='11987654321';password='qa-password-123'}
$user2 = Request 'register recipient' POST '/api/v1/auth/register' 201 @{name='QA Recipient';cpf='98765432100';email=$email2;phone='11987654322';password='qa-password-123'}
# The API has no PIX-key registration endpoint; seed only our isolated QA database.
if ($DatabaseContainer -notmatch '^wallet-qa-' -or $DatabaseName -notmatch '^qa_[a-z_]+$') { throw 'Use an isolated wallet-qa- container and qa_ database.' }
$recipientID = [guid]$user2.data.id
$sql = "INSERT INTO pix_keys (user_id, key, type, status) VALUES ('$recipientID', '$email2', 'email', 'active');"
docker exec $DatabaseContainer psql -U qa_user -d $DatabaseName -v ON_ERROR_STOP=1 -c $sql | Out-Null
if ($LASTEXITCODE -ne 0) { throw 'Failed to seed the isolated PIX key.' }
$null = Request 'duplicate registration rejected' POST '/api/v1/auth/register' 400 @{name='Duplicate';cpf='12345678909';email=$email1;phone='11987654321';password='qa-password-123'}
$null = Request 'wrong password rejected' POST '/api/v1/auth/login' 401 @{email=$email1;password='incorrect'}
$login1 = Request 'sender login' POST '/api/v1/auth/login' 200 @{email=$email1;password='qa-password-123'}
$login2 = Request 'recipient login' POST '/api/v1/auth/login' 200 @{email=$email2;password='qa-password-123'}
$token1 = $login1.data.access_token
$token2 = $login2.data.access_token
$null = Request 'profile' GET '/api/v1/auth/me' 200 $null $token1
AssertBalance $token1 0
$null = Request 'credit' POST '/api/v1/wallet/credit' 200 @{amount=100;description='QA credit'} $token1
$null = Request 'negative credit rejected' POST '/api/v1/wallet/credit' 400 @{amount=-1} $token1
$null = Request 'debit' POST '/api/v1/wallet/debit' 200 @{amount=10;description='QA debit'} $token1
$null = Request 'insufficient debit rejected' POST '/api/v1/wallet/debit' 400 @{amount=1000} $token1
$null = Request 'transfer' POST '/api/v1/transactions/transfer' 200 @{to_user_id=$user2.data.id;amount=20;description='QA transfer'} $token1
$null = Request 'self transfer rejected' POST '/api/v1/transactions/transfer' 400 @{to_user_id=$user1.data.id;amount=1} $token1
$null = Request 'PIX' POST '/api/v1/pix' 200 @{pix_key=$email2;amount=5} $token1
$null = Request 'unknown PIX key rejected' POST '/api/v1/pix' 400 @{pix_key='missing@example.com';amount=1} $token1
$null = Request 'boleto' POST '/api/v1/boleto' 200 @{amount=15} $token1
$null = Request 'zero boleto rejected' POST '/api/v1/boleto' 400 @{amount=0} $token1
$null = Request 'transaction history' GET '/api/v1/transactions?limit=10&offset=0' 200 $null $token1
AssertBalance $token1 65
AssertBalance $token2 25
$qaToken = $token1
$qaBaseUrl = $BaseUrl
$statuses = 1..10 | ForEach-Object -Parallel {
    $response = Invoke-WebRequest -Uri "$using:qaBaseUrl/api/v1/wallet/credit" -Method Post -ContentType 'application/json' -Body '{"amount":1,"description":"QA concurrent credit"}' -Headers @{Authorization="Bearer $using:qaToken"} -SkipHttpErrorCheck -TimeoutSec 15
    [int]$response.StatusCode
} -ThrottleLimit 10
if (@($statuses | Where-Object { $_ -ne 200 }).Count) { throw 'Concurrent credit failed.' }
$script:checks += 10
AssertBalance $token1 75
Write-Host 'PASS: 10 concurrent credits preserved the balance.'
$statuses = 1..10 | ForEach-Object -Parallel {
    $response = Invoke-WebRequest -Uri "$using:qaBaseUrl/api/v1/wallet/debit" -Method Post -ContentType 'application/json' -Body '{"amount":10,"description":"QA concurrent debit"}' -Headers @{Authorization="Bearer $using:qaToken"} -SkipHttpErrorCheck -TimeoutSec 15
    [int]$response.StatusCode
} -ThrottleLimit 10
if (@($statuses | Where-Object { $_ -eq 200 }).Count -ne 7 -or @($statuses | Where-Object { $_ -eq 400 }).Count -ne 3) { throw 'Concurrent debit did not enforce available balance.' }
$script:checks += 10
AssertBalance $token1 5
Write-Host 'PASS: concurrent debits rejected overdrafts.'
$refresh = Request 'refresh token' POST '/api/v1/auth/refresh' 200 @{refresh_token=$login1.data.refresh_token}
$null = Request 'refreshed access token' GET '/api/v1/auth/me' 200 $null $refresh.data.access_token
$null = Request 'invalid refresh rejected' POST '/api/v1/auth/refresh' 401 @{refresh_token='invalid'}
Write-Host "Completed: $script:checks checks passed."
