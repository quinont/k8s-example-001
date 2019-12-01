FROM python:3.8.0-alpine

EXPOSE 80
RUN mkdir app
WORKDIR /app
COPY . .
RUN pip install -r requirements.txt
CMD ["python", "main.py"]
