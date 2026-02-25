# Szablony odpowiedzi (triage)

Gotowe szablony do wklejania w Issues i Discussions podczas tygodnia feedbacku.
Kopiuj treść, dostosuj `[...]` do konkretnego przypadku.

---

## 1. needs-repro

```
Dziękujemy za zgłoszenie!

Żeby móc potwierdzić błąd, potrzebujemy minimalnych kroków reprodukcji:

1. Dokładna wersja docflow (`./build/docflow --version`)
2. OS + architektura
3. Kolejność komend (od zera, w czystym katalogu)
4. Plik `docflow.yaml` lub minimalny przykład konfiguracji
5. Pełne wyjście z terminala (usuń sekrety/tokeny jeśli są)

Dodaję label `status/needs-repro`. Gdy pojawią się kroki, wrócimy do oceny.
Jeśli nie otrzymamy odpowiedzi w ciągu 7 dni, zgłoszenie zostanie automatycznie zamknięte.
```

**Labele do ustawienia:** `status/needs-repro`, `type/bug`

---

## 2. needs-info

```
Dziękujemy za zgłoszenie!

Żeby ocenić ten problem, potrzebujemy kilku dodatkowych informacji:

- [ ] [opisz brakującą informację 1]
- [ ] [opisz brakującą informację 2]

Dodaję label `status/needs-info`. Po uzupełnieniu wrócimy do oceny.
```

**Labele do ustawienia:** `status/needs-info`

---

## 3. use-discussions (pytanie w Issues)

```
Dziękujemy za kontakt!

To pytanie/pomysł najlepiej omówić w GitHub Discussions, gdzie możemy prowadzić
wątkową rozmowę bez zaśmiecania kolejki Issues:

https://github.com/Aerendal/doc-flow/discussions

Przenoszę tam tę rozmowę i zamykam to Issue. Do zobaczenia w Discussions!
```

**Akcja:** zamknij Issue przez "Convert to Discussion" (GitHub UI) lub zamknij z tym komentarzem.

---

## 4. duplicate

```
To zgłoszenie jest duplikatem #[numer] ([tytuł]).

Dodaję label `status/duplicate` i zamykam to Issue. Dalsza dyskusja w #[numer].
```

**Labele:** `status/duplicate`  
**Akcja:** zamknij Issue.

---

## 5. wontfix (z uzasadnieniem)

```
Dziękujemy za zgłoszenie!

Po ocenie zdecydowaliśmy, że ta zmiana jest poza zakresem projektu / nie pasuje
do aktualnej roadmapy z następujących powodów:

> [krótkie uzasadnienie]

Jeśli masz inne podejście lub dodatkowy kontekst, chętnie porozmawiamy
w GitHub Discussions: https://github.com/Aerendal/doc-flow/discussions
```

**Labele:** `status/wontfix`  
**Akcja:** zamknij Issue.
