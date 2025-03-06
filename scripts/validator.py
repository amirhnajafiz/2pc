transactions = []
balances = {}

with open("inputs.txt", "r") as file:
    for line in file:
        x = line.replace("(", "").replace(")", "").split(", ")
        transactions.append({
            "sender": x[0],
            "receiver": x[1],
            "amount": int(x[2])
        })
        if x[0] not in balances:
            balances[x[0]] = 10
        if x[1] not in balances:
            balances[x[1]] = 10

with open("output.txt", "w") as file:
    for trx in transactions:
        file.write(f"({trx['sender']}, {trx['receiver']}, {trx['amount']}) => ")
        if balances[trx['sender']] >= trx['amount']:
            balances[trx['sender']] -= trx['amount']
            balances[trx['receiver']] += trx['amount']
            file.write("OK")
        else:
            file.write("Failed")
        file.write(f"\t{trx['sender']}:{balances[trx['sender']]}, {trx['receiver']}:{balances[trx['receiver']]}\n")
