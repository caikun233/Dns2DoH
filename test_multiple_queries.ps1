# Test multiple DNS query types
Write-Host "Testing multiple DNS queries..." -ForegroundColor Cyan

# Query 1: A record
nslookup -type=A baidu.com 127.0.0.1 2>&1 | Out-Null
Start-Sleep -Milliseconds 500

# Query 2: AAAA record  
nslookup -type=AAAA google.com 127.0.0.1 2>&1 | Out-Null
Start-Sleep -Milliseconds 500

# Query 3: MX record
nslookup -type=MX gmail.com 127.0.0.1 2>&1 | Out-Null
Start-Sleep -Milliseconds 500

Write-Host "`nWaiting for logs to be written..." -ForegroundColor Yellow
Start-Sleep -Seconds 2

Write-Host "`nLast 4 query logs:" -ForegroundColor Green
Write-Host "=" * 80 -ForegroundColor Gray

Get-Content logs\queries.log | Select-Object -Last 4 | ForEach-Object {
    $log = $_ | ConvertFrom-Json
    Write-Host "`nQuery: $($log.domain) ($($log.query_type))" -ForegroundColor Cyan
    Write-Host "  Client: $($log.client_ip)" -ForegroundColor Gray
    Write-Host "  DoH Server: $($log.doh_server)" -ForegroundColor Gray
    Write-Host "  Response: $($log.response_code) | Duration: $($log.duration_ms)ms" -ForegroundColor Gray
    
    if ($log.answers -and $log.answers.Count -gt 0) {
        Write-Host "  Answers:" -ForegroundColor Yellow
        foreach ($answer in $log.answers) {
            Write-Host "    [$($answer.type)] $($answer.name) -> $($answer.value) (TTL: $($answer.ttl))" -ForegroundColor White
        }
    } else {
        Write-Host "  No answers" -ForegroundColor Gray
    }
}

Write-Host "`n" + ("=" * 80) -ForegroundColor Gray
