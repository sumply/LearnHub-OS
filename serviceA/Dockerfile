FROM archlinux:latest

RUN pacman -Syu --noconfirm && \
	pacman -S --noconfirm base-devel cmake boost gcc clang git && \
	pacman -Scc --noconfirm

WORKDIR /app
COPY . .
ARG DEBUG_FLAG=""
RUN cmake -DCMAKE_EXPORT_COMPILE_COMMANDS=ON -DCMAKE_CXX_FLAGS="${DEBUG_FLAG}" -B build
RUN cmake --build build
ENTRYPOINT ["./build/serviceA"]
