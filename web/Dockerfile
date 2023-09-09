FROM node:18 as builder
RUN npm install -g pnpm

WORKDIR /app

COPY .npmrc package.json pnpm-lock.yaml ./
RUN pnpm install

COPY . .
RUN pnpm build

FROM node:18
RUN npm install -g pnpm

WORKDIR /app

COPY --from=builder /app/package.json /app/pnpm-lock.yaml ./
RUN pnpm install 
COPY --from=builder /app/build ./

EXPOSE 3000
CMD ["node", "./index.js"]