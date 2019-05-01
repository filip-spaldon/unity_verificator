int tmp;
int N;
int a[5];
init {
N = 5;
a[0] = 85;
a[1] = 12;
a[2] = 33;
a[3] = 25;
a[4] = 23;
}
active proctype process_0() {
do
:: a[0] > a[1] ->
atomic {
tmp = a[0];
a[0] = a[1];
a[1] = tmp;
}
:: else -> skip
od
}
active proctype process_1() {
do
:: a[1] > a[2] ->
atomic {
tmp = a[1];
a[1] = a[2];
a[2] = tmp;
}
:: else -> skip
od
}
active proctype process_2() {
do
:: a[2] > a[3] ->
atomic {
tmp = a[2];
a[2] = a[3];
a[3] = tmp;
}
:: else -> skip
od
}
active proctype process_3() {
do
:: a[3] > a[4] ->
atomic {
tmp = a[3];
a[3] = a[4];
a[4] = tmp;
}
:: else -> skip
od
}
