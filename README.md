## Run development

- Build storekeeper
```shell
./gradlew :storekeeper:jibDockerBuild
```
- Build tristar provider
```shell
./gradlew :providers:tristar:jibDockerBuild
```
- Export environment variables from `.env.dev` for current context
```shell
export $(xargs < .env.dev)
```
- Run docker-compose
```shell
docker-compose up -f docker-compose.dev.yml up
```