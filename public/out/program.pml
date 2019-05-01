int tmp;
int x;
int y;
init {
x = 60;
y = 48;
}
active proctype process_0() {
do
:: x > y ->
atomic {
x  =  x - y
}
:: else -> skip
od
}
active proctype process_1() {
do
:: y > x ->
atomic {
y  =  y - x
}
:: else -> skip
od
}
