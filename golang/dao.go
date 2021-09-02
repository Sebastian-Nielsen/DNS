func getMoney(c Caller, amount int) {
   if (amount <= account[c].amount) { 
      sendMoney(c, amount);
      account[c].amount = account[c].amount - amount;
   }
}
