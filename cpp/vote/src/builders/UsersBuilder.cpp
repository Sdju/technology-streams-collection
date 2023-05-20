//
// Created by zede on 20.05.2023.
//

#include "UsersBuilder.h"

[[nodiscard]] std::shared_ptr<User> UsersBuilder::create(const std::string& name, const std::string& pwd, int age) const {
  return std::make_shared<User>(name, pwd, age);
}