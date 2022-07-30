#syntax=docker/dockerfile:1

FROM golang:1.18-alpine AS builder
LABEL maintainer="svetlankasvetka@gmail.com"

#створюємо робочу директорію
WORKDIR /app

#скопіювали go файл
COPY *.go ./
#створили файл, що виконується
RUN go build -o gses2

#створюємо ще образ , де просто запускаємо цей файл. Можна було б використовувати scratch замість alpine
FROM alpine
COPY --from=builder /app/gses2 /app/gses2
#скопіювали txt файл
COPY *.txt ./app

#які порти потрібні - 3333 - на ньому програма слухає, 2525 - відправлення листів, 443 - https запит 
EXPOSE 3333 2525 443

#запускаємо 
CMD ["app/gses2"]
#або
#ENTRYPOINT ["app/gses2"]