cmake_minimum_required(VERSION 3.28)
project(serviceA)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

add_executable(${PROJECT_NAME}
    src/listener.cxx
    src/main.cxx)

target_include_directories(${PROJECT_NAME} PRIVATE ${CMAKE_SOURCE_DIR}/include)

find_package(Boost REQUIRED)
include(FetchContent)
FetchContent_Declare(
  nlohmann_json
  GIT_REPOSITORY https://github.com/nlohmann/json.git
  GIT_TAG master
)
FetchContent_MakeAvailable(nlohmann_json)
FetchContent_Declare(
	jwt_cpp
	GIT_REPOSITORY https://github.com/Thalhammer/jwt-cpp.git
	GIT_TAG master
)
FetchContent_MakeAvailable(jwt_cpp)

target_link_libraries(${PROJECT_NAME}
    PRIVATE Boost::headers
    PRIVATE jwt-cpp::jwt-cpp
    PRIVATE nlohmann_json::nlohmann_json
)
