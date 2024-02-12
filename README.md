# pipeline-demo-golang

Esse projeto configura uma pipeline de CI/CD usando GitHub Actions em uma aplicação Golang para fazer as seguintes tarefas:

- Criação e publicação de imagem docker no DockerHub 
- Verificação de linter no código.
- Verificação de testes unitários e de integração.
- Verificação de build
- Deploy no servioço Google Cloud Run.

Você precisa seguir alguns passos.

## Gerar imagem e enviar ao Dockerhub ao comentar

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

Aqui está um exemplo de um arquivo de workflow que constrói e publica uma imagem Docker no Docker Hub quando um comentário "dockerhub" é feito em um pull request:

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

## Realizar teste ao comentar

O arquivo test-on-comment.yml é um arquivo de workflow do GitHub Actions configurado para executar testes unitários em resposta a um comentário específico feito em um Pull Request. A estrutura do arquivo é definida em partes distintas para estabelecer gatilhos, trabalhos e etapas a serem executadas.

### Início do arquivo: Definição do nome e evento de gatilho

```yaml
name: PR Comment Trigger

on:
  issue_comment:
    types:
      - created
```

**Descrição**:
- `name: PR Comment Trigger`: Define o nome do workflow como "PR Comment Trigger".
- `on: issue_comment: types: [- created]`: Especifica que esse workflow é ativado quando um comentário é criado numa issue. Como Pull Requests (PRs) são tratados como tipos especiais de issues no GitHub, isso significa que o workflow será ativado por comentários em Pull Requests também.

### Job: `comment-trigger`

```yaml
jobs:
  comment-trigger:
    runs-on: ubuntu-latest
```

**Descrição**:
- Define um job chamado `comment-trigger`.
- `runs-on: ubuntu-latest`: Especifica que o job será executado no ambiente Ubuntu mais recente disponível.

### Passos para executar:

#### Checkout do Código

```yaml
- name: Check out code
  uses: actions/checkout@v2
```

**Descrição**:
- `name: Check out code`: Define o nome desta etapa como "Check out code".
- `uses: actions/checkout@v2`: Usa a action `checkout@v2` para clonar o repositório do código-fonte para o ambiente de execução do GitHub Actions, possibilitando a realização de testes no código.

#### Execução Condicionada por Comentário no PR

```yaml
- name: Run on PR comment
  if: github.event_name == 'issue_comment' && contains(github.event.comment.body, 'teste')
  run: |
    echo "The trigger comment was found in the PR comment. Running your job now."
    go test -cover -race -vet=off ./...
```

**Descrição**:
- `name: Run on PR comment`: Define o nome desta etapa como "Run on PR comment".
- `if: github.event_name == 'issue_comment' && contains(github.event.comment.body, 'teste')`: Este if condicional verifica duas coisas: primeiro, se o evento que desencadeou o workflow foi um comentário em issue (que inclui PRs); segundo, se o corpo do comentário contém a palavra "teste". Somente se ambas as condições forem verdadeiras o workflow prossegue para a execução dos comandos dentro do bloco `run`.
- O primeiro comando `echo` dentro do bloco `run` é uma mensagem de confirmação de que o comentário de gatilho foi encontrado e o job será executado.
- O segundo comando `go test -cover -race -vet=off ./...` executa os testes unitários do código Go, com flags que habilitam a verificação de cobertura de código, a detecção de condições de corrida e desabilita o 'vet' (uma ferramenta de análise estática de código Go) para todos os subpacotes do diretório atual.

## Realizar deploy no Cloud Run ao comentar

O arquivo deploy-on-comment é um arquivo de workflow do GitHub Actions configurado para realizar o deploy de uma aplicação no Google Cloud Run quando um comentário contendo a palavra "deploy" é feito em um Pull Request (PR).

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