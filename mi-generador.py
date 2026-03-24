import sys

def generate_output_file(output_filename, amount_clients):
    compose = f"""name: tp0
services:
  server:
    container_name: server
    image: server:latest
    volumes: 
      - ./server/config.ini:/config.ini
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
    networks:
      - testing_net
"""
    for i in range(1, amount_clients + 1):
        compose += f"""
  client{i}:
    container_name: client{i}
    image: client:latest
    volumes: 
      - ./client/config.yaml:/config.yaml
      - ./.data/agency-{i}.csv:/agency.csv
    entrypoint: /client
    environment:
      - CLI_ID={i}
      - CLI_NOMBRE=Santiago Lionel
      - CLI_APELLIDO=Lorca
      - CLI_DOCUMENTO=30904465
      - CLI_NACIMIENTO=1999-03-17
      - CLI_NUMERO=7574
    networks:
      - testing_net
    depends_on:
      - server
"""
    
    compose += """
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""

    with open(output_filename, "w") as f:
        f.write(compose)
    
    print(f"Archivo de salida con el nombre {output_filename} generado con {amount_clients} cantidad de clientes.")

if __name__ == "__main__":
    output_filename = sys.argv[1]
    amount_clients = int(sys.argv[2])
    
    generate_output_file(output_filename, amount_clients)