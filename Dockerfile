#syntax=docker/dockerfile:1

FROM golang:1.18-alpine AS builder
LABEL maintainer="svetlankasvetka@gmail.com"

#��������� ������ ���������
WORKDIR /app

#��������� go ����
COPY *.go ./
#�������� ����, �� ����������
RUN go build -o gses2

#��������� �� ����� , �� ������ ��������� ��� ����. ����� ���� � ��������������� scratch ������ alpine
FROM alpine
COPY --from=builder /app/gses2 /app/gses2
#��������� txt ����
COPY *.txt ./app

#�� ����� ������ - 3333 - �� ����� �������� �����, 2525 - ����������� �����, 443 - https ����� 
EXPOSE 3333 2525 443

#��������� 
CMD ["app/gses2"]
#���
#ENTRYPOINT ["app/gses2"]