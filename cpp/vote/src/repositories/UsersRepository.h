//
// Created by zede on 20.05.2023.
//

#ifndef VOTE_SRC_REPOSITORIES_USERSREPOSITORY_H_
#define VOTE_SRC_REPOSITORIES_USERSREPOSITORY_H_

#include "../entities/User.h"

class UsersRepository {
 public:
  User* addUser(std::shared_ptr<User>&& user);

  [[nodiscard]] std::shared_ptr<User> findUser (const std::string& name, const std::string& pwd) const;

  void logTable();

 private:
  std::unordered_map<std::string, std::shared_ptr<User>> users_;
};


#endif //VOTE_SRC_REPOSITORIES_USERSREPOSITORY_H_
