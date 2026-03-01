import subprocess

IMAGE_NAME="learnhub-server"
TAG_NAME="sumply/learnhub-server"
VERSION="latest"

def build_image() -> void:
    run(f"docker build . -t {IMAGE_NAME}:{VERSION}")

def tag_image() -> void:
    run(f"docker tag {IMAGE_NAME}:{VERSION} {TAG_NAME}:{VERSION}")

def push_image() -> void:
    run(f"docker push {TAG_NAME}")

def login() -> void:
    run(f"docker login")

def run(cmd: str) -> void:
    print(f"running {cmd}")
    subprocess.run(cmd, shell=True, check=True)

if __name__ == "__main__":
    login()
    build_image()
    tag_image()
    push_image()