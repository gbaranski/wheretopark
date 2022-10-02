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


# Webapp

- Build shared library
```
./gradlew :shared:jsNodeProductionLibraryDistribution
```
- Change directory to webapp/
```
cd webapp
```
- Edit `.env.local` file with required environment variables defined in `next.config.js`
```
- Install dependencies
```
yarn install
```
- Run development server
```
yarn dev
```