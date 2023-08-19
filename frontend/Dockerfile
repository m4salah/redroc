FROM node:lts AS build
WORKDIR /app

COPY . .

RUN npm install
RUN npm run build

FROM node:lts-alpine as runtime

WORKDIR /app

COPY --from=build /app/dist ./dist
COPY --from=build /app/node_modules ./node_modules

ENV HOST=0.0.0.0
ENV PORT=3000
EXPOSE 3000

CMD node ./dist/server/entry.mjs