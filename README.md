# exemplo-curso-docker

Comando para compilar a imagem Docker:

```bash
cd exemplo-curso-docker
docker build -t pehrolim/exemplo-curso-docker:1.0.0 .
```

```bash
docker container run -d \
-p 85:3000 \
--rm \
--name=exemplo \
phenriquerolim/exemplo-curso-docker:1.0.0
```

```bash
docker login -u CHANGE_HERE
 
# Crie o repositório CHANGE_HERE/exemplo-curso-docker no Docker Hub antes de enviar a imagem
 
docker push CHANGE_HERE/exemplo-curso-docker:1.0.0
```

