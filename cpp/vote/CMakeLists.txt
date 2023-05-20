cmake_minimum_required(VERSION 3.25)
project(vote)

set(CMAKE_CXX_STANDARD 20)

add_executable(vote src/entities/user.cpp src/entities/user.h src/entities/session.cpp src/entities/session.h src/builders/SessionBuilder.cpp src/builders/SessionBuilder.h src/builders/UsersBuilder.cpp src/builders/UsersBuilder.h src/main.cpp src/utilies/utilities.h src/utilies/utilities.cpp src/repositories/UsersRepository.cpp src/repositories/UsersRepository.h src/repositories/SessionsRepository.cpp src/repositories/SessionsRepository.h src/entities/Cookie.cpp src/entities/Cookie.h src/controllers/BaseController.cpp src/controllers/BaseController.h src/entities/ServerData.cpp src/entities/ServerData.h src/controllers/HomeController.cpp src/controllers/HomeController.h src/controllers/AuthController.cpp src/controllers/AuthController.h)

find_package(unofficial-http-parser REQUIRED)
find_package(fmt REQUIRED)
find_package(restinio REQUIRED)

target_link_libraries(vote PRIVATE restinio::restinio)
