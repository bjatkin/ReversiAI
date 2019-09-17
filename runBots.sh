START_SERVER="cd ./ReversiJava/ReversiServer && javac *.java && java Reversi 30 "
START_JAVA_AI="cd ./ReversiJava/ReversiRandom_Java && javac *.java && java RandomGuy 127.0.0.1 1 "
START_GO_AI="go build && ./ReversiBot 127.0.0.1 2 "
WAIT="sleep 5 "

exec $WAIT && $START_GO_AI
eval $START_SERVER