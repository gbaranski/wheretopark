FROM node:18 as builder
RUN curl -f https://get.pnpm.io/v6.16.js | node - add --global pnpm

WORKDIR /app

COPY .npmrc package.json pnpm-lock.yaml ./
RUN pnpm install

COPY . .
RUN pnpm build

FROM node:18
RUN curl -f https://get.pnpm.io/v6.16.js | node - add --global pnpm

WORKDIR /app

COPY --from=builder /app/package.json /app/pnpm-lock.yaml ./
RUN pnpm install 
COPY --from=builder /app/build ./

EXPOSE 3000
CMD ["node", "./index.js"]