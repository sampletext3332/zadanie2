## Sprawozdanie Zadanie 2
### autor: Mateusz Ł. 101619

### Kroki łańcucha

#### Przygotowanie środowiska
Na początku workflow pobierany jest kod źródłowy z repozytorium za pomocą actions/checkout@v4.<br>
Następnie inicjalizowane są narzędzia potrzebne do budowania obrazów Docker, w tym:
```docker/setup-qemu-action``` – umożliwia budowanie obrazów dla wielu architektur (linux/amd64, linux/arm64),
```docker/setup-buildx-action``` – aktywuje Buildx, czyli rozszerzony mechanizm budowania obrazów.

#### Logowanie do GH, Dockerhub
DockerHub – przy użyciu DOCKERHUB_USERNAME oraz DOCKERHUB_TOKEN<br>

#### Budowanie obrazu do testów (scan)
Tworzony jest obraz Docker tylko dla architektury linux/amd64.<br>
Obraz ten nie jest publikowany, jest ładowany lokalnie (load: true),<br>
otrzymuje tag:
```local-scan:<github_sha>```

#### Skanowanie 
Skan Trivy obejmuje podatności HIGH i CRITICAL
<br>jeśli wykryte zostaną krytyczne problemy pipeline zostanie przerwany
<br>jeśli skan przejdzie pomyślnie, obraz zostanie publikowany

#### Budowanie i publikacja obrazu produkcyjnego
W przypadku braku wykrycia podatności HIGH/CRITICAL obraz zostanie opublikowany na GHCR


### Tagowanie
W workflow zastosowano automatyczne tagowanie przez ```docker/metadata-action```<br>
```flavor: latest=false``` - brak tagu latest

### Używane zmienne/secrets

```DOCKERHUB_TOKEN``` - token do logowania do DockerHub<br>
```CACHE_REPO``` - repozytorium używane do cache buildów Dockera<br>
```CACHE_TAG``` - tag dla cache<br>
```DOCKERHUB_USERNAME``` - login konta DockerHub<br>

#### link do Dockerhub
[https://hub.docker.com/r/sampletext333/zadanie2-cache](https://hub.docker.com/r/sampletext333/zadanie2-cache)
