# DNS 查询测试脚本
# 使用原始 UDP socket 向 dns2doh 发送 DNS 查询

$dnsServer = "127.0.0.1"
$dnsPort = 15353
$domain = "google.com"

Write-Host "正在查询 $domain ..." -ForegroundColor Cyan

try {
    # 创建 UDP 客户端
    $udpClient = New-Object System.Net.Sockets.UdpClient
    $udpClient.Connect($dnsServer, $dnsPort)
    
    # 构建 DNS 查询包（简单的 A 记录查询）
    # DNS 查询格式: Transaction ID(2) + Flags(2) + Questions(2) + Answer RRs(2) + Authority RRs(2) + Additional RRs(2) + Query
    $transactionId = Get-Random -Maximum 65535
    $query = [byte[]]@(
        # Header
        [byte]($transactionId -shr 8), [byte]($transactionId -band 0xFF),  # Transaction ID
        0x01, 0x00,  # Flags: Standard query
        0x00, 0x01,  # Questions: 1
        0x00, 0x00,  # Answer RRs: 0
        0x00, 0x00,  # Authority RRs: 0
        0x00, 0x00   # Additional RRs: 0
    )
    
    # 添加域名查询部分
    foreach ($label in $domain.Split('.')) {
        $query += [byte]$label.Length
        $query += [System.Text.Encoding]::ASCII.GetBytes($label)
    }
    $query += 0x00  # 结束符
    
    # Query type (A) and class (IN)
    $query += 0x00, 0x01  # Type: A
    $query += 0x00, 0x01  # Class: IN
    
    # 发送查询
    $bytesSent = $udpClient.Send($query, $query.Length)
    Write-Host "已发送 $bytesSent 字节查询" -ForegroundColor Yellow
    
    # 设置超时
    $udpClient.Client.ReceiveTimeout = 5000
    
    # 接收响应
    $remoteEP = New-Object System.Net.IPEndPoint([System.Net.IPAddress]::Any, 0)
    $response = $udpClient.Receive([ref]$remoteEP)
    
    Write-Host "收到 $($response.Length) 字节响应" -ForegroundColor Green
    Write-Host "`n响应数据（十六进制）:" -ForegroundColor Cyan
    
    # 显示前 100 字节的十六进制数据
    $hexString = ""
    for ($i = 0; $i -lt [Math]::Min(100, $response.Length); $i++) {
        $hexString += "{0:X2} " -f $response[$i]
        if (($i + 1) % 16 -eq 0) {
            Write-Host $hexString
            $hexString = ""
        }
    }
    if ($hexString) { Write-Host $hexString }
    
    # 简单解析响应获取 IP 地址
    $answerCount = ($response[6] -shl 8) -bor $response[7]
    Write-Host "`n应答记录数: $answerCount" -ForegroundColor Cyan
    
    if ($answerCount -gt 0) {
        Write-Host "✓ DNS 查询成功！" -ForegroundColor Green
        Write-Host "服务器正常工作，可以解析 DNS 请求。" -ForegroundColor Green
    }
    
    $udpClient.Close()
    
} catch {
    Write-Host "错误: $_" -ForegroundColor Red
} finally {
    if ($udpClient) {
        $udpClient.Close()
    }
}
