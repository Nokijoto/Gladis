# Discord Bot z Google Gemini AI

Ten dokument opisuje kroki potrzebne do stworzenia, optymalizacji i wdrożenia bota na Discordzie, który działa jako proxy dla API Google Gemini. Projekt zostanie napisany w języku Go, z naciskiem na minimalne zużycie zasobów (RAM, CPU).

## Plan Działania:

#### Faza 1: Inicjalizacja i Konfiguracja Projektu

1.  **Inicjalizacja Repozytorium Git**: Stworzenie lokalnego repozytorium do śledzenia zmian w kodzie (`git init`).
2.  **Stworzenie Pliku `.gitignore`**: Zdefiniowanie plików, które Git ma ignorować (np. skompilowane binarki).
3.  **Inicjalizacja Modułu Go**: Uruchomienie `go mod init` w celu zarządzania zależnościami projektu.
4.  **Stworzenie Głównej Struktury Projektu**: Utworzenie pliku `main.go` jako punktu wejściowego aplikacji.

#### Faza 2: Implementacja Rdzenia Bota

5.  **Dodanie Zależności**: Pobranie niezbędnych bibliotek: `discordgo` do obsługi Discorda i oficjalnego SDK Google dla Gemini.
6.  **Konfiguracja Aplikacji**: Implementacja wczytywania kluczy API (Discord Bot Token, Gemini API Key) ze zmiennych środowiskowych, aby uniknąć zapisywania ich w kodzie. Stworzenie pliku `.env.example` jako szablonu.
7.  **Podstawowa Logika Bota Discord**: Nawiązanie połączenia z API Discorda i implementacja prostej komendy testowej (np. `!ping`), aby potwierdzić, że bot działa.

#### Faza 3: Integracja z Google Gemini

8.  **Implementacja Klienta Gemini**: Stworzenie funkcji, która będzie wysyłać zapytania tekstowe do API Google Gemini i odbierać odpowiedzi.
9.  **Połączenie Logiki Bota z Gemini**: Zintegrowanie klienta Gemini z handlerem wiadomości na Discordzie. Bot będzie przechwytywał wiadomości, przesyłał je do AI i odsyłał odpowiedź na kanał.

#### Faza 4: Optymalizacja i Wdrożenie

10. **Obsługa Błędów i Logowanie**: Zapewnienie stabilności przez dodanie obsługi błędów (np. problemów z siecią) i prostego logowania zdarzeń.
11. **Optymalizacja Zużycia Pamięci**: Przegląd kodu w celu zminimalizowania alokacji pamięci i zapewnienia płynnego działania na 256 MB RAM.
12. **Przygotowanie do Wdrożenia (Dockerfile)**: Stworzenie `Dockerfile` z wykorzystaniem "multi-stage build", aby skompilować aplikację do małej, statycznej binarki i umieścić ją w minimalnym obrazie kontenera (np. `scratch` lub `alpine`).
13. **Dokumentacja Projektu**: Zaktualizowanie pliku `README.md` o instrukcje dotyczące budowania, konfiguracji i uruchamiania bota.

## Instalacja Go:

Ponieważ Go nie jest zainstalowany, oto kroki, które możesz wykonać, aby go zainstalować:

1.  **Pobierz Go**: Odwiedź oficjalną stronę [https://go.dev/dl/](https://go.dev/dl/) i pobierz najnowszą stabilną wersję dla swojego systemu operacyjnego.
2.  **Zainstaluj Go**: Postępuj zgodnie z instrukcjami instalacji dla Twojego systemu. Zazwyczaj obejmuje to rozpakowanie archiwum i dodanie katalogu `bin` Go do zmiennej środowiskowej `PATH`.
    *   **Linux/macOS**:
        ```bash
        # Przykład dla Linuxa, ścieżka może się różnić
        wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz # Zmień wersję i architekturę jeśli potrzebujesz
        sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc # Lub ~/.bashrc jeśli używasz bash
        source ~/.zshrc # Lub source ~/.bashrc
        ```
    *   **Windows**: Uruchom pobrany instalator MSI.
3.  **Zweryfikuj Instalację**: Otwórz nowy terminal i wpisz `go version`. Powinieneś zobaczyć zainstalowaną wersję Go.

Po zainstalowaniu Go, będziemy mogli kontynuować z inicjalizacją modułu Go (`go mod init gladis`).
