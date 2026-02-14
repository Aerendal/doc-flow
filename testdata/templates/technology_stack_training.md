---
title: Technology Stack Training
status: needs_content
---

# Technology Stack Training

## Metadane
- Właściciel: [Engineering Enablement/L&D]
- Wersja: v0.2
- Data aktualizacji: RRRR-MM-DD
- Status: draft | in review | approved

## Cel dokumentu
Program szkolenia z firmowego stacku technologicznego (frontend/backend/data/cloud). Ma skrócić ramp‑up, ujednolicić praktyki i zwiększyć jakość delivery.

## Zakres i granice
- Obejmuje: przegląd architektury, języki/frameworki, standardy kodu/testów/CI-CD, observability, bezpieczeństwo (IAM/secret/PII), dane (DB, cache, messaging), IaC/deploy, narzędzia deweloperskie, A11y/UX (frontend), procedury on-call.  
- Poza zakresem: ogólne kursy CS (algorytmy itp.).

## Wejścia i wyjścia
- Wejścia: opis architektury, standardy kodu, runbooki, narzędzia, przykładowe projekty, wymagania bezpieczeństwa, syllabus L&D.  
- Wyjścia: sylabus i materiały (slajdy, laby), checklisty DoR/DoD, lab repos, oceny (quiz/praktyka), harmonogram i lista uczestników, wyniki ewaluacji.

## Powiązania (meta)
- Key Documents: coding_guidelines, observability_plan, security_requirements, ci_cd_standards, architecture_vision, onboarding_engineer, on_call_training.  
- Key Document Structures: architektura, kod/testy, deploy/IaC, security, data, tools, ewaluacja.  
- Document Dependencies: repos, CI/CD, cloud accounts/sandboxes, monitoring, LMS.

## Zależności dokumentu
Wymaga: aktualnych standardów kodu/CI-CD/security, dostępu do repo i sandboxes, materiałów architektonicznych, listy trenerów i uczestników, LMS. Braki = DoR otwarte.

## Powiązania sekcja↔sekcja
- Architektura → Kod/standardy → Deploy/IaC → Observability/on-call.  
- Security → CI/CD → Deploy/ops.  
- Laby → Oceny → Ewaluacja/udoskonalenia.

## Fazy cyklu życia
- Przygotowanie materiałów i labów.  
- Szkolenia i oceny.  
- Ewaluacja i iteracje.  
- Odświeżenia cykliczne przy zmianie stacku.

## Struktura sekcji
1) Cele szkolenia i grupa docelowa  
2) Architektura i przegląd systemu  
3) Języki/frameworki i standardy kodu/testów  
4) CI/CD i IaC (pipelines, deployments, feature flags)  
5) Observability i on-call (metryki/logi/traces, alerty, runbooki)  
6) Bezpieczeństwo (IAM, secrets, PII, dependency scanning)  
7) Dane (DB, cache, queues, schema migration)  
8) Narzędzia dev (IDE, linters, formatters, local env)  
9) Laby, oceny i certyfikacja  
10) Ryzyka, decyzje, otwarte pytania

## Wymagane rozwinięcia
- Sylabus modułów i czas; laby per domena.  
- Checklista środowiska lokalnego i dostępów.  
- Quiz/praktyka i rubryka ocen.  
- Plan aktualizacji materiałów (release train).

## Wymagane streszczenia
- One‑pager: cele, moduły, terminy, prerekwizyty.  
- Snapshot wyników: frekwencja, zdawalność, NPS.

## Guidance (skrót)
- Ucz na realnych repo i case’ach; automatyzuj laby.  
- Wbuduj security/observability w każdy moduł.  
- Sprawdzaj prerekwizyty i dostęp do narzędzi przed startem.  
- Zbieraj feedback po każdej sesji; iteruj materiały.  
- Odświeżaj przy zmianie stacku (releases).

## Szybkie powiązania
- linkage_index.jsonl (technology/stack/training)  
- coding_guidelines, observability_plan, security_requirements, ci_cd_standards

## Jak używać dokumentu
1. Przygotuj sylabus/laby i dostęp do środowisk.  
2. Przeprowadź szkolenia; oceń wiedzę; zbierz feedback.  
3. Aktualizuj materiały i linkage_index; odnotuj DoR/DoD.

## Checklisty Definition of Ready (DoR)
- [ ] Standardy kodu/CI-CD/security aktualne.  
- [ ] Repo/laby i sandboxes gotowe.  
- [ ] Trenerzy i terminy potwierdzeni; LMS skonfigurowany.  
- [ ] Prerekwizyty/instalacje zakomunikowane.  
- [ ] Rubryka ocen przygotowana.

## Checklisty Definition of Done (DoD)
- [ ] Sesje przeprowadzone; frekwencja/wyniki zapisane; status/wersja/data uzupełnione.  
- [ ] Materiały/laby zaktualizowane wg feedbacku.  
- [ ] Certyfikacje/badges wydane (jeśli dotyczy); linkage_index uzupełniony.  
- [ ] Plan kolejnego odświeżenia ustalony.  
- [ ] Ryzyka/lessons learned zapisane.

## Definicje robocze
- IaC: Infrastructure as Code.  
- NPS: Net Promoter Score.  
- Sandboxes: izolowane środowiska treningowe.

## Przykłady użycia
- Onboarding nowych inżynierów.  
- Migracja stacku (np. monolit → microservices).  
- Program upskilling dla zespołów legacy.

## Ryzyka i ograniczenia
- Brak dostępu do narzędzi → słaba efektywność.  
- Nieaktualne materiały → złe praktyki.  
- Niska frekwencja → niewielki efekt.

## Decyzje i uzasadnienia
- Format (self‑paced vs live).  
- Zakres labów vs czas.  
- Kryteria zaliczenia i certyfikacji.

## Założenia
- Dostępne środowiska i narzędzia.  
- Budżet na trenerów/czas.  
- Zespoły mają czas na udział.

## Otwarte pytania
- Jak mierzyć wpływ na velocity/quality?  
- Jakie moduły obowiązkowe vs opcjonalne?  
- Jak często aktualizować sylabus?

## Powiązania z innymi dokumentami
- onboarding_engineer — plan startowy.  
- on_call_training — operacje.  
- security_requirements — bezpieczeństwo.

## Wymagane odwołania do standardów
- Wewnętrzne standardy kodu/CI/CD/security/observability.  
- Polityki PII/RODO jeśli obejmuje dane.
