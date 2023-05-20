//
// Created by zede on 20.05.2023.
//

#include "UsersRepository.h"

User* UsersRepository::addUser(std::shared_ptr<User>&& user) {
  auto result = users_.try_emplace(user->name, std::move(user));
  if (result.second) {
    return result.first->second.get();
  }
  return nullptr;
}

[[nodiscard]] std::shared_ptr<User> UsersRepository::findUser (const std::string& name, const std::string& pwd) const {
  auto result = users_.find(name);
  if (result == users_.end()) {
    return nullptr;
  }

  if (result->second->pwd != pwd) {
    return nullptr;
  }
  return result->second;
}

void UsersRepository::logTable() {
  fmt::print("USERS_TABLE:\n");
  for (auto& user : users_) {
    fmt::print("{}\n", *(user.second));
  }
}