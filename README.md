# pipeline-demo-golang

Esse projeto configura uma pipeline de CI/CD usando GitHub Actions em uma aplicação Golang para fazer as seguintes tarefas:

- Verificação da integração contínua e qualidade do código
  - Verificação de linter no código.
  - Verificação de testes unitários e de integração.
  - Verificação de build
- Criação e publicação de imagem docker no DockerHub 
- Deploy no servioço Google Cloud Run.

## Verificação da integração contínua e qualidade do código

O arquivo continous-integration.yaml define um pipeline de integração contínua (CI) nomeado "CI pipeline" que executa várias ações relacionadas à verificação de conformidade com a nomenclatura de branches, linting, teste e build de código, orientado principalmente para projetos Golang. Abaixo, você encontrará uma explicação detalhada de cada etapa desse pipeline:

### Gatilhos do Workflow

- **on:** Este pipeline é acionado por dois eventos principais:
    - **push**: Quando há um push na branch `master` (incluindo tags, embora a especificação de tags esteja vazia aqui, indicando que o push de tags não ativa este pipeline).
    - **pull_request**: Em qualquer pull request que vise a branch `master`.

### Jobs (Tarefas)

#### 1. `check_gitflow_conformance`
- **Name**: GitFlow Branch Naming
- **runs-on**: Executa em um runner com Ubuntu.
- **steps**: Contém a etapa para verificação do nome da branch.
    - Verifica se o nome da branch segue as convenções do GitFlow (`feature/`, `bugfix/`, `hotfix/`, `release/`). Se não seguir, o script sinaliza erro e encerra o processo com `exit 1`.

#### 2. `golangci`
- **Name**: Lint
- **runs-on**: Executa em um runner com Ubuntu.
- **steps**: Efetua a configuração da versão do Go, checa o código e executa o linter no código Go usando a action `golangci/golangci-lint-action@v3`.
    - **Set up Go**: Configura o ambiente com Go versão 1.18.
    - **Check out code**: Faz checkout do código.
    - **Lint Go Code**: Executa o golangci-lint, que é uma ferramenta de análise estática para o código Go.

#### 3. `test`
- **Name**: Test
- **runs-on**: Executa em um runner com Ubuntu.
- **steps**: Prepara o ambiente Go, checa o código e executa testes unitários.
    - **Set up Go**: Configura o Go 1.18.
    - **Check out code**: Faz checkout do código.
    - **Run unit Tests.**: Executa testes unitários com cobertura e opção `-race` para detecção de race conditions, excluindo verificações com `go vet`.

#### 4. `build`
- **Name**: Build
- **runs-on**: Executa em um runner com Ubuntu.
    - **needs**: Esta tarefa depende da conclusão bem-sucedida das tarefas `golangci`, `test` e `check_gitflow_conformance`.
- **steps**: Prepara Go, checa o código e compila o projeto.
    - **Set up Go**: Configura Go 1.18.
    - **Check out code**: Faz checkout do código.
    - **Build**: Compila o projeto usando o comando `go build`.

Cada job é executado em uma instância separada do Ubuntu na versão mais recente disponível para GitHub Actions, garantindo que o código seja revisado por padrões de nomenclatura, qualidade (via linting), por testes automatizados e, finalmente, que seja compilável. A dependência explícita do job `build` para os demais jobs (`golangci`, `test`, `check_gitflow_conformance`) assegura que a compilação só acontecerá se todas as verificações preliminares tiverem sucesso.

## Criação e publicação de imagem docker no DockerHub

### Pré-requisitos:

1. **Conta no Docker Hub:** Você precisa de uma conta no Docker Hub para publicar imagens. Se não tiver, crie uma em [Docker Hub](https://hub.docker.com/).

2. **Repositório GitHub:** Você deve ter um repositório no GitHub com acesso para configurar Actions.

3. **Dockerfile:** O projeto no GitHub deve ter um Dockerfile. Assim, o GitHub Actions pode construir a imagem Docker baseada neste arquivo.

### Passos para Configuração:

#### 1. Configurar Secrets do GitHub

Primeiro, configure os segredos (secrets) no repositório GitHub para armazenar suas credenciais do Docker Hub de forma segura.

- Acesse o repositório GitHub > Settings > Secrets > New repository secret.
- Adicione dois segredos:
  - `DOCKER_USERNAME`: seu nome de usuário do Docker Hub.
  - `DOCKER_PASSWORD`: sua senha ou token de acesso do Docker Hub.

#### 2. Criando o Workflow do GitHub Actions

- No seu repositório, navegue até a aba "Actions" > "New workflow" > "set up a workflow yourself" ou crie um novo arquivo `.yml` em `.github/workflows` no seu repositório. Por exemplo: `.github/workflows/docker-publish.yml`.

Aqui está o código do arquivo de workflow que constrói e publica uma imagem Docker no Docker Hub quando um comentário "dockerhub" é feito em um pull request:

```yaml
name: Build and Publish Docker image

on:
  issue_comment:
    types: [created]

jobs:
  build-and-publish:
    if: github.event.issue.pull_request && contains(github.event.comment.body, 'dockerhub')
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v2

      - name: Login to DockerHub
        uses: docker/login-action@v1
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Build and push Docker image
        uses: docker/build-push-action@v2
        with:
          push: true
          tags: seuusuario/nomeimagem:latest

```

**Notas:**
- Substitua `seuusuario/nomeimagem` pelo seu usuário do Docker Hub e pelo nome que deseja dar à sua imagem.
- Este script especificamente procura por comentários em pull requests (não em issues comuns, que também são tecnicamente issues no GitHub). Se o comentário contém a palavra "dockerhub", então o job de construir e publicar a imagem é ativado.

#### 3. Acionando o Workflow

Para acionar este workflow, vá até um pull request existente no seu repositório e adicione um comentário com o texto "dockerhub". O GitHub Actions verificará automaticamente o comentário, e se corresponder à condição definida (`contains(github.event.comment.body, 'dockerhub')`), iniciará o processo de construção e publicação da imagem no Docker Hub.

## Realizar deploy no Cloud Run ao comentar

O arquivo clou-run-deploy-on-comment é um arquivo de workflow do GitHub Actions configurado para realizar o deploy de uma aplicação no Google Cloud Run quando um comentário contendo a palavra "deploy" é feito em um Pull Request (PR).

### Fluxo de funcionamento

Durante o processo de build, o Cloud Build é configurado para compilar o código e construir a imagem do contêiner. Após a conclusão do build, O Cloud Build utiliza o Artifact Registry para armazenar as imagens de contêineres e outros artefatos em um repositório. O repositório do Artifact Registry por sua vez utiliza o Cloud Storage por debaixo dos panos para armazenar os dados. Imagens de contêiner podem ser construídas e armazenadas no Cloud Storage antes de serem referenciadas na ação para deploy.


### Pré-Requisitos de funcionamento

1. **Google Cloud Project**: Certifique-se de que você tenha um projeto no Google Cloud.
2. **Enable Cloud Run API & Google Cloud Build API**: Certifique-se de que as APIs do Cloud Run e do Cloud Build estejam habilitadas no seu projeto Google Cloud.
3. **Service Account**: Crie uma Service Account no Google Cloud com permissões suficientes para fazer deploy no Cloud Run e configure uma chave de acesso. Guarde o JSON da chave de acesso com segurança.

4. **Configure Secrets no GitHub**: No seu repositório GitHub, vá até `Settings > Secrets` e adicione os seguintes secrets:
   - `GCP_SERVICE_NAME`: O nome do serviço onde será feito o deploy (esse serviço já deve estar criado antes da pipeline ser executada)
   - `GCP_PROJECT_ID`: O ID do seu projeto no Google Cloud.
   - `GCP_SA_KEY`: O conteúdo do arquivo JSON da sua chave de acesso da Service Account (codificado em Base64).
   - `GCP_REGION`: A região onde o serviço será implantado, por exemplo: us-central1

### Analise do código

#### Configuração de autenticação do Google Cloud

```yaml
- name: Setup Google Cloud Authentication
  uses: google-github-actions/auth@v2
  with:
    credentials_json: ${{ secrets.GCP_SA_KEY }}
```

**Descrição**:
- Utiliza a action `google-github-actions/auth@v2` para configurar a autenticação no Google Cloud usando uma chave de conta de serviço (Service Account Key), que é armazenada de forma segura no GitHub Secrets (`GCP_SA_KEY`).

#### Configuração do projeto do Google Cloud

```yaml
- name: Configure Google Cloud Project
  run: gcloud config set project ${{ secrets.GCP_PROJECT_ID }}
```

**Descrição**:
- Configura o gcloud, a CLI do Google Cloud, para usar o ID do projeto especificado na secreta `GCP_PROJECT_ID`.

#### Fazer o deploy no Cloud Run

```yaml
- name: Deploy to Cloud Run
  id: deploy
  uses: google-github-actions/deploy-cloudrun@v2
  with:
    service: ${{ secrets.GCP_SERVICE_NAME }}
    region: ${{ env.GCP_REGION }}
    source: ./
```

**Descrição**:
- Utiliza a action `google-github-actions/deploy-cloudrun@v2` para fazer o deploy do código no serviço especificado na secret `GCP_SERVICE_NAME`, na região especificada na variável de ambiente `GCP_REGION`, tendo como fonte o diretório atual do repositório.

#### Mostrar Output

```yaml
- name: Show Output
  run: echo ${{ steps.deploy.outputs.url }}
```

**Descrição**:
- Mostra a URL resultante do deploy no Cloud Run, facilitando a verificação do resultado do deploy.