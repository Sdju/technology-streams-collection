//
// Created by zede on 20.05.2023.
//

#include "utilities.h"

std::string generate_uuid() {
  std::random_device dev;
  std::mt19937 rng(dev());
  std::uniform_int_distribution<std::mt19937::result_type> distHex(0,15);

  auto getChar = [&distHex, &rng](){
    auto val = distHex(rng);
    if (val < 10) {
      return static_cast<char>(val + '0');
    }
    return static_cast<char>((val - 10) + 'A');
  };

  std::stringstream uuid_stream;
  for (int i = 0; i < 8; i += 1) {
    uuid_stream << getChar();
  }
  uuid_stream << '-';
  for (int i = 0; i < 4; i += 1) {
    uuid_stream << getChar();
  }
  uuid_stream << '-';
  for (int i = 0; i < 4; i += 1) {
    uuid_stream << getChar();
  }
  uuid_stream << '-';
  for (int i = 0; i < 4; i += 1) {
    uuid_stream << getChar();
  }
  uuid_stream << '-';
  for (int i = 0; i < 12; i += 1) {
    uuid_stream << getChar();
  }

  return uuid_stream.str();
}