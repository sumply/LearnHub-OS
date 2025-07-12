FROM archlinux:latest AS builder

RUN pacman -Syu --noconfirm && \
	pacman -S --noconfirm base-devel cmake boost gcc clang git && \
	pacman -Scc --noconfirm

WORKDIR /app
COPY . .

ARG DEBUG_FLAG=""

RUN cmake -DCMAKE_EXPORT_COMPILE_COMMANDS=ON -DCMAKE_CXX_FLAGS="${DEBUG_FLAG}" -B build
RUN cmake --build build

FROM archlinux:latest

WORKDIR /app

COPY --from=builder /app/build/serviceA ./serviceA

ENTRYPOINT ["./serviceA"]
