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