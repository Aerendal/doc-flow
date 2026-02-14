# Security Verification Guide

Ten dokument opisuje trzy poziomy weryfikacji release:
- `baseline`: checksums
- `strong`: checksums + cosign
- `enterprise`: checksums + cosign + attestation + SBOM

## Baseline (checksums)
Pobierz z tej samej strony release:
- archiwum dla swojej platformy (`.tar.gz` lub `.zip`)
- `checksums.txt`

Linux:
```bash
sha256sum -c checksums.txt
```

macOS (fallback):
```bash
shasum -a 256 -c checksums.txt
```

Windows (PowerShell, single file):
```powershell
$zip="docflow-windows-amd64.zip"
$expected = (Select-String -Path .\checksums.txt -Pattern $zip).ToString().Split(" ")[0]
$actual = (Get-FileHash .\$zip -Algorithm SHA256).Hash.ToLower()
if ($expected.ToLower() -ne $actual) { throw "SHA256 mismatch" } else { "OK" }
```

## Strong (checksums + cosign)
Po poprawnym checksumie zweryfikuj podpisy:
- `<archive>.sig` i `<archive>.cert`
- `checksums.txt.sig` i `checksums.txt.cert`

Archiwum:
```bash
cosign verify-blob \
  --certificate docflow-linux-amd64.tar.gz.cert \
  --signature   docflow-linux-amd64.tar.gz.sig \
  docflow-linux-amd64.tar.gz
```

Checksums:
```bash
cosign verify-blob \
  --certificate checksums.txt.cert \
  --signature   checksums.txt.sig \
  checksums.txt
```

## Enterprise (checksums + cosign + attestation + SBOM)
Oprócz kroków `strong`:
- sprawdź w GitHub Release/Actions, że dostępne są `Attestations` dla artefaktów
- pobierz i przejrzyj `sbom.cdx.json` (CycloneDX)

To jest tryb do audytów supply-chain i compliance.

## FAQ
`checksums.txt` nie pasuje do nazw plików:
Upewnij się, że pliki mają oryginalne nazwy z release (bez sufiksów typu `(1)`).

Brak `*.sig`/`*.cert`:
Release może być starszy lub podpisy mogły nie zostać wygenerowane. Wtedy użyj co najmniej poziomu `baseline`.

Czy da się weryfikować offline:
Checksumy tak. `cosign verify-blob` z lokalnym blob/sig/cert też działa lokalnie, ale pełna walidacja łańcucha zewnętrznego może wymagać sieci.

macOS nie ma `sha256sum`:
Użyj `shasum -a 256 -c checksums.txt`.
