C_YELLOW='\[\e[01;33m\]'
C_RED='\[\e[00;31m\]'
C_MAGENTA='\[\e[01;35m\]'
C_LB_BLUE='\[\e[01;36m\]'
C_RESET='\[\e[0m\]'
C_LB_GREEN='\[\e[01;32m\]'
C_GOLD='\[\e[00;33m\]'
C_GREEN='\[\e[00;32m\]'
C_D_BLUE='\[\e[01;34m\]'
C_GRAY='\[\e[00;90m\]'
C_B_RED='\[\e[01;91m\]'
C_BG_GRAY='\[\e[00;44m\]'

LIGHTGREEN="\033[1;32m"
LIGHTBLUE="\033[1;34m"
LIGHTYELLOWBLINK="\033[5;33m"
LIGHTYELLOW="\033[1;33m"
LIGHTRED="\033[1;31m"
LIGHMAGENTA="\033[1;95m"
parse_git_branch() {
    #git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/ (\1)/'
    git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/\1/'
}
check_git_branch() {
	IS_GIT=`git status 2> /dev/null | wc -l`
	if [ $IS_GIT -eq 0 ]; then
		exit
	fi
    GIT_DIFF=`git diff --name-only 2> /dev/null | wc -l`
    GIT_DIFF_CACHED=`git diff --cached --name-only 2> /dev/null | wc -l`
    if [ $GIT_DIFF -eq 0 ] && [ $GIT_DIFF_CACHED -eq 0 ]; then
		GIT_UP_TO_DATE=`git status 2> /dev/null | grep 'up-to-date\|up to date' | wc -l`
		if [ $GIT_UP_TO_DATE -eq 0 ]; then
			#echo -e -n "$LIGHMAGENTA**";
			echo -n "**"
		else
			#echo -e -n "$LIGHTBLUE+";
			echo -n "+"
		fi
    elif [ $GIT_DIFF -eq 0 ]; then
        #echo -e -n "$LIGHTYELLOW*";
	echo "*"
    else
        #echo -e -n "$LIGHTRED*";
	echo "*"
    fi
}
function getExitCode(){
    VAR=$?;
    if [ "$VAR" -eq 0 ]; then
        #echo -e -n "$LIGHTGREEN";
	echo -n ""
    else
        #echo -e -n "$LIGHTYELLOWBLINK";
	echo -n ""
    fi
    echo $VAR;
}


if [ "$color_prompt" = yes ]; then
    #PS1="$C_BG_GRAY $C_YELLOW\$(getExitCode)$C_BG_GRAY $C_RESET$C_RED\u$C_LB_BLUE@$C_RED\h$C_LB_GREEN:$C_GOLD\w$C_GRAY [$C_GREEN\$(parse_git_branch)$C_LB_GREEN\$(check_git_branch)$C_GRAY] $C_D_BLUE-> $C_RESET"
    PS1="$C_BG_GRAY $C_YELLOW\$(getExitCode)$C_BG_GRAY $C_RESET${C_RED}p$C_LB_BLUE@${C_RED}wsl$C_LB_GREEN:$C_GOLD\w$C_GRAY [$C_GREEN\$(parse_git_branch)$C_LB_GREEN\$(check_git_branch)$C_GRAY] $C_D_BLUE-> $C_RESET"
else
    PS1='${debian_chroot:+($debian_chroot)}\u@\h:\w\$ '
fi

RED='\033[01;31m'
searchInFiles(){
    if [ $# -eq 0 ]; then
        echo -e "${RED}not enough parameters"
    elif [ $# -eq 1 ]; then
        echo -e "${RED}searching..."
        grep --include=\*.go -rni . -e "${1}"  2>/dev/null
    elif [ $# -eq 2 ]; then
        grep --include=\*."${2}" -rni . -e "${1}" 2>/dev/null
    else
        echo -e "${RED}too much parameters"
    fi
}
alias s='searchInFiles'

alias repo="cd /mnt/d/repo/goApps/src/coreCommunicationFramework"
alias difff="git diff  -- . ':(exclude)*.spam' -- . ':(exclude)*.db'"

export GOROOT=/usr/local/go
export GOPATH=/mnt/d/repo/goApps
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

export L=pawel33317
export P=pass