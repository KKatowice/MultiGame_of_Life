import kaplay from "kaplay";
import "kaplay/global";
//import WebSocket from 'ws';

const k = kaplay({
    width: window.width,
    height: window.height,
    font: "sans-serif",
    canvas: document.querySelector("#mycanvas"),
    background: [ 0, 0, 0, ],
})
let ws;

//consts
const SPEED = 360;
let userId




scene('startGame',()=>{
    k.setBackground([255,0,0])
    /* add([
        text("Create Game", { width: width() / 3 }),
        pos(12, 12),
    ]); */

    /* add([
        text("Choose Name:"),
        pos(0,100)
    ]) */
    let objName = add([
        text("Choose Name: "),
        textInput(),
        pos(0,0)
    ])
    /* onkeydown = (event) => {
        console.log(`?????? ${event.key} ${event.key[0]}`);
        
        objName.text += event.key[0]
    }; */

    let  objConfirm = add([
        rect(250, 50),
        color(255, 255, 255),
        pos(0,100),
        area(),
        "buttonConf"
    ])
    let  txtConfirm = add([
        text("Confirm"),
        color(0, 0, 0),
        pos(50,110),
    ])
    onClick("buttonConf",()=>{
        let res = fetch(`http://127.0.0.1:3000/api/temp_user?nome=${objName.text.split(" Name: ")[1]}`).then(result => {
            
           let jsn = result.text()
           jsn.then(resultz => {
                console.log(`res: ${resultz}`);
                userId = resultz
           })
           if(result.ok){
            go('game')
           }
        })
        
        
        
    })
    
})


scene('game',async()=>{
    //
    const queryString = window.location.search
    const urlParams = new URLSearchParams(queryString);
    let rmId = urlParams.get("roomId")
    if (!rmId){
        
        ///genera
        let res = fetch(`http://127.0.0.1:3000/api/host_game?h=${windowHeight}&w=${windowWidth}&userId=${userId}`).then(result => {
            
            let jsn = result.text()
            jsn.then(resultz => {
                 console.log(`res: ${resultz}`);
                 rmId = resultz
                 window.location.href = window.location.href + "/" + roomId ;
            })
         })
    }else{
        //joina
        let res = fetch(`http://127.0.0.1:3000/api/join_game?userId=${userId}&roomID=${rmId}`).then(result => {
            let jsn = result.text()
            jsn.then(resultz => {
                 console.log(`res: ${resultz}`);
                 let bho = resultz
            })
         })
    }
    ws = new window.WebSocket(`ws://127.0.0.1:3001/ws/${rmId}`);
    //per wss

    k.setBackground([0,0,0])

    // Add player game object (a white square)
    const player = add([
        rect(50, 50), // Create a rectangle (width, height)
        color(255, 255, 255), // Set color to white (RGB values)
        pos(center()), // Position it at the center of the screen
    ]);

    // onKeyDown() registers an event that runs every frame as long as the user is holding a certain key
    onKeyDown("a", () => {
        // Move player left
        if(player.pos.x>=0){
            player.move(-SPEED, 0);
        }
    });
    onKeyDown("d", () => {
        // Move player right
        
        if(player.pos.x+player.width <= k.width()){
            player.move(SPEED, 0);
        }
    });
    onKeyDown("w", () => {
        // Move player up
        if(player.pos.y>=0){
            player.move(0, -SPEED);
        }
        
    });
    onKeyDown("s", () => {
        // Move player down
        if(player.pos.y+player.height <= k.height()){
            player.move(0, SPEED);
        }
        
    });
    onKeyDown(["w","a","s","d"],()=>{
        
        ws.send({'uid':'1','pos':player.pos})
    })


    onClick(() => {
        console.log(mousePos())
    })


    add([
        text("Press arrow to move", { width: width() / 3 }),
        pos(12, 12),
    ]);
    add([
        text("Press .  to stop", { width: width() / 3 }),
        pos(350, 15),
    ]); 

    let enemInit = 30
    let enemySize = 15; 
    let isEnemInit = false
    let spacing = 10; // Spaziatura tra gli enemies
    let windowWidth = window.width();
    let windowHeight = window.height();

    function createEnemy(posX,posY){
        return [
            rect(enemySize, enemySize),
            pos(posX,posY),
            color(255, 0, 0),
            "enemy",
        ]
    }


    function calculateGrid(windowWidth, windowHeight, enemySize, spacing) {
        // Calcola quante colonne possono stare nella larghezza della finestra
        let cols = Math.floor(windowWidth / (enemySize + spacing));
        
        // Calcola quante righe possono stare nell'altezza della finestra
        let rows = Math.floor(windowHeight / (enemySize + spacing));
        
        return { rows, cols };
    }

    function calculateGridPositions(grid, windowWidth, windowHeight, enemySize, spacing) {
        let paddingX = (windowWidth - (grid.cols * enemySize)) / (grid.cols + 1); // Spaziatura orizzontale
        let paddingY = (windowHeight - (grid.rows * enemySize)) / (grid.rows + 1); // Spaziatura verticale
        
        let positions = [];
        for (let r = 1; r <= grid.rows; r++) {
            for (let c = 1; c <= grid.cols; c++) {
                let x = c * (enemySize + paddingX);
                let y = r * (enemySize + paddingY);
                positions.push({ x: x, y: y });
            }
        }
        return positions;
    }

    function test_spawnPosition(gp) {
        let newP = [];
        
        let availablePositions = [...gp];
        
        let enemyCount = 0;
        if (isEnemInit) {
            while (enemyCount < enemInit && availablePositions.length > 0) {
                let randomIndex = getRandomInt(availablePositions.length); 
                let randomPosition = availablePositions.splice(randomIndex, 1)[0]; // Rimuovi la posizione scelta
                let e = createEnemy(randomPosition.y, randomPosition.x); // Usa la posizione casuale
                newP.push(e);
                enemyCount++;
            }
            //qua mando posizione iniziale a server
            isEnemInit = true
        } else {
            //qua ricezione msg wss per 
            while (enemyCount < enemInit && availablePositions.length > 0) {
                let randomIndex = getRandomInt(availablePositions.length); 
                let randomPosition = availablePositions.splice(randomIndex, 1)[0]; 
                let e = createEnemy(randomPosition.y, randomPosition.x); 
                newP.push(e);
                enemyCount++;
            }
        }

        return newP;
    }


    const sleep = (delay) => new Promise((resolve) => setTimeout(resolve, delay))


    //main
    let cycles = true
    let grid = calculateGrid(windowWidth, windowHeight, enemySize, spacing);
    let gridPositions = calculateGridPositions(grid, windowWidth, windowHeight, enemySize, spacing);
    onKeyDown(".", () => {
        cycles = false
    });
    while (cycles){
        ///console.log("siiiiiiiiiiiii??")
        let tst = test_spawnPosition(gridPositions)
        for (let i = 0; i < tst.length; i++) {
            add(tst[i])
            //console.log("si??")
        }
        
        await sleep(1000)
        destroyAll("enemy")
        
        await sleep(10)
    }
	
})

go("startGame")
//go("game")








function getRandomInt(max) {
	return Math.floor(Math.random() * max);
  }