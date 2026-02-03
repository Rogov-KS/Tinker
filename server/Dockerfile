FROM node:18-alpine AS deps
WORKDIR /app

COPY package*.json tsconfig*.json ./
COPY prisma ./prisma
RUN npm ci

FROM deps AS builder
WORKDIR /app
COPY tsconfig*.json ./
COPY src ./src
COPY prisma ./prisma
RUN npx prisma generate
RUN npm run build

FROM node:18-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

# Install openssl for Prisma
RUN apk add --no-cache openssl

COPY --from=deps /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist
COPY package*.json ./
COPY prisma ./prisma
COPY entrypoint.sh ./

RUN chmod +x entrypoint.sh

EXPOSE 3000

ENTRYPOINT ["./entrypoint.sh"]
