$allbctl = "C:\allbctl-test\allbctl.exe"
$pass = 0
$fail = 0

function Run-Test {
    param($label, $testArgs)
    Write-Host ""
    Write-Host "=== $label ===" -ForegroundColor Cyan
    $output = & $allbctl @testArgs 2>&1
    $code = $LASTEXITCODE
    Write-Host $output
    if ($code -ne 0) {
        Write-Host "FAIL (exit $code)" -ForegroundColor Red
        $script:fail++
    } else {
        Write-Host "PASS" -ForegroundColor Green
        $script:pass++
    }
}

Run-Test "version"          @("version")
Run-Test "help"             @("--help")
Run-Test "status (full)"    @("status")
Run-Test "status projects"  @("status", "projects")
Run-Test "status runtimes"  @("status", "runtimes")
Run-Test "bootstrap status" @("bootstrap", "status")

Write-Host ""
Write-Host "========================================"
Write-Host "Results: $pass passed, $fail failed"
if ($fail -gt 0) {
    Write-Host "OVERALL: FAIL" -ForegroundColor Red
    exit 1
} else {
    Write-Host "OVERALL: PASS" -ForegroundColor Green
}
