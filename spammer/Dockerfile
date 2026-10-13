FROM python:3.13-slim

WORKDIR /app

# No requirements.txt needed: the script only uses Python's standard library
COPY spammer.py .

CMD ["python", "-u", "spammer.py"]