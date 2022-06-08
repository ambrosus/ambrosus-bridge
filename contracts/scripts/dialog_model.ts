import inquirer from "inquirer";

class Dialog {
    output = (data: string) => console.log(data);

    confirmation = async(message: string) => {
        const {confirmation} = await inquirer.prompt([
            {
                type: 'list',
                name: 'confirmation',
                message: message,
                choices: [
                    {
                        name: "NO",
                        value: false
                    },
                    {
                        name: "YES",
                        value: true
                    }
                ]
            }
        ]);
        return confirmation;
    }

    askToChooseFromArray = async (arr: Array<any>, message: string) => {
        const {value} = await inquirer.prompt([
            {
                type: 'list',
                name: 'value',
                message: message,
                choices: arr
            }
        ]);
        return value;
    }
}

export default new Dialog();