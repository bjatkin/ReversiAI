START_SERVER="cd ./ReversiJava/ReversiServer && javac *.java && java Reversi 30"
START_JAVA_AI="cd ./ReversiJava/ReversiRandom_Java && javac *.java && java RandomGuy 127.0.0.1 1"
START_GO_AI="cd go build && ./ReversiBot 127.0.0.1 2"

exec $START_SERVER
sleep 5
exec $START_JAVA_AI
exec $START_GO_AI