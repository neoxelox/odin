# #!/usr/bin/env bash

# RED=$(tput setaf 203)
# GREEN=$(tput setaf 46)
# YELLOW=$(tput setaf 226)
# BLUE=$(tput setaf 39)
# RESET=$(tput sgr0)

# RELEASE=odin
# NAMESPACE=asgard

# if ! command -v helm &> /dev/null
# then
#     echo "${RED}* helm could not be found, install it from https://helm.sh/docs/intro/install/${RESET}"
#     exit 1
# fi

# echo "${BLUE}* helm version${RESET}"
# helm version
# echo ""

# echo "${BLUE}* List revision history for $RELEASE ${RESET}"
# helm history $RELEASE -n $NAMESPACE
# echo ""

# if [ "$1" == "--run" ]; then
#     echo "${YELLOW}-------------------------------------------------${RESET}"
#     echo "${YELLOW}WARNING:${RESET} Rollback does not affect migrations."
#     echo "${YELLOW}-------------------------------------------------${RESET}"

#     read -p "I undertand the risk and want to continue [Y/n] " -n 1
#     echo    # move to a new line
#     if [[ $REPLY =~ ^[Nn]$ ]]; then
#         echo "${RED}Skipped.${RESET}"
#         exit 1
#     fi

#     echo "${GREEN}* Running rollback${RESET}"
#     helm rollback $RELEASE -n $NAMESPACE
#     if [ $? -ne 0 ]; then
#         echo "${RED}* Rollback failed ❌ ${RESET}"
#         helm history $RELEASE -n $NAMESPACE
#         exit 1
#     else
#         helm history $RELEASE -n $NAMESPACE
#     fi
# else
#     echo "${BLUE}* Running in ${RED}dry mode${BLUE}, rollback not executed, execute it with ${YELLOW}--run${BLUE} to perform the rollback${RESET}"
# fi
