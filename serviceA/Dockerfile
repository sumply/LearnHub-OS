FROM archlinux:latest

RUN pacman -Syu --noconfirm && \
	pacman -S --noconfirm base-devel cmake boost gcc clang git && \
	pacman -Scc --noconfirm

WORKDIR /app
COPY . .
RUN cmake -DCMAKE_EXPORT_COMPILE_COMMANDS=ON -B build
RUN cmake --build build
ENTRYPOINT ["./build/serviceA"]
