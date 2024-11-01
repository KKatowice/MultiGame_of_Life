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
//let userId




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
        let userId
        let res = fetch(`http://127.0.0.1:3000/api/temp_user?nome=${objName.text.split(" Name: ")[1]}`).then(result => {
            
           let jsn = result.text()
           jsn.then(resultz => {
                console.log(`res: ${resultz}`);
                userId = resultz
                if(result.ok){
                    go('game',userId)
                }
           })
           
        })
        
        
        
    })
    
})


scene('game',async(uid)=>{
    console.log(`USERID??? ${uid}`);
    
    let windowWidth;
    let windowHeight;
    const urlz = window.location.href.split("/")
    console.log(`urlszz`, urlz);
    let rmId
    let uids = {}

    

    if (urlz[3] == ''){
        windowWidth = window.width()
        windowHeight = window.height()
        console.log(`fra1???/`,windowWidth,windowHeight);
        
        ///genera
        let result = await fetch(`http://127.0.0.1:3000/api/host_game?h=${windowHeight}&w=${windowWidth}&userId=${uid}`)
            
        let resultz = await result.text()
            
        console.log(`ressesss: ${resultz}`);
        rmId = resultz
        const currentUrl = window.location.href;
        const newUrl = currentUrl + rmId;
        history.pushState(null, '', newUrl);
            
         
    }else{
        //joina
        rmId = urlz[3]
        let result = await fetch(`http://127.0.0.1:3000/api/join_game?userId=${uid}&roomID=${rmId}`)
        let resultz = await result.json()
    
        ////console.log(`resseesd: ${resultz}`);
        let wh = resultz
        windowWidth = wh.w
        windowHeight = wh.h
        
         ////console.log(`fra2???/`,windowWidth,windowHeight);
    }
    
    ws = new window.WebSocket(`ws://127.0.0.1:3001/ws/${rmId}`);
    
    //per wss
    ws.addEventListener("message", (event) => { ///TODO MMH MANNAGGIO MESA PROBLEMA A INVVIATUTTO INSIEME SU UN CHANNEL?/ BHO E SU JS ? BHO 
        let jsEvent = JSON.parse(event.data)
        console.log(`jsevent`, jsEvent);
        
        if(jsEvent.uid){
            if (parseInt(jsEvent.uid) != uid){
                //console.log("Message from other player ", event.data, "io->",uid);
                if (!uids[jsEvent.uid]){
                    const mplayer = add([
                        rect(40, 40), // Create a rectangle (width, height)
                        color(255, 255, 50), // Set color to white (RGB values)
                        pos(jsEvent.x,jsEvent.y), // Position it at the center of the screen
                    ])
                    uids[jsEvent.uid] = mplayer
                }
                uids[jsEvent.uid].moveTo(vec2(jsEvent.x,jsEvent.y),SPEED)
            }
        }else if (jsEvent.RID){
            destroyAll("enemy")
            jsEvent.Positions.forEach(element => {
                add(createEnemy(element.X, element.Y)); // Usa la posizione casuale
                
            });
            ////console.log(`------------------------`);
            
        }
        
      });
      
        ws.onclose = function (e){
            let result = fetch(`http://127.0.0.1:3000/api/quit_lobby?userId=${uid}`)
        };
        window.addEventListener("unload", function () {
            if(ws.readyState == WebSocket.OPEN)
                ws.close();
        });

      function createEnemy(posX,posY){

        return [
            rect(15, 15), ///todo
            pos(posX,posY),
            color(255, 0, 0),
            "enemy",
        ]
    }
    k.setBackground([0,0,0])

    // Add player game object (a white square)
    const player = add([
        rect(50, 50), // Create a rectangle (width, height)
        color(255, 255, 255), // Set color to white (RGB values)
        pos(0,0), // Position it at the center of the screen
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
        //console.log(`pos? send to ws`, player.pos);
        ws.send(JSON.stringify({'uid':parseInt(uid),'x':parseFloat(player.pos.x), 'y':parseFloat(player.pos.y)}))
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

    


    /* function calculateGrid(windowWidth, windowHeight, enemySize, spacing) {
        console.log(`wwwidjiwjd`,windowWidth, windowHeight);
        
        // Calcola quante colonne possono stare nella larghezza della finestra
        let cols = Math.floor(windowWidth / (enemySize + spacing));
        
        // Calcola quante righe possono stare nell'altezza della finestra
        let rows = Math.floor(windowHeight / (enemySize + spacing));
        console.log(`rwo ${rows} -- col ${cols}`);
        
        return { rows, cols };
    } */

/*     function calculateGridPositions(grid, windowWidth, windowHeight, enemySize, spacing) {
        let paddingX = (windowWidth - (grid.cols * enemySize)) / (grid.cols + 1); // Spaziatura orizzontale
        let paddingY = (windowHeight - (grid.rows * enemySize)) / (grid.rows + 1); // Spaziatura verticale
        //console.log(`@@@@@@@`,grid, windowWidth, windowHeight, enemySize, spacing);
        
        let positions = [];
        let rowsList = []///tst
        for (let r = 1; r <= grid.rows; r++) {
            let trl = []
            for (let c = 1; c <= grid.cols; c++) {
                let x = c * (enemySize + paddingX);
                let y = r * (enemySize + paddingY);
                positions.push({ x: x, y: y });
                trl.push({ x: x, y: y })///tst
                console.log(`{${x}:${y}}`);
                
            }
            rowsList.push(trl)///tst
        }
        
        
        return positions;
    } */

   /*  function test_spawnPosition(gp) {
        let newP = [];
        
        gp.forEach(element => {
            let e = createEnemy(element.y, element.x); // Usa la posizione casuale
            console.log(`enem?`,e);
                    
            add(e) 
        });
        
        
        
        return newP;
    } */


    ///const sleep = (delay) => new Promise((resolve) => setTimeout(resolve, delay))

    /* console.log(`????????????????????''''''''`,windowWidth, windowHeight);
    
    //main
    let cycles = true
    //let grid = calculateGrid(windowWidth, windowHeight, enemySize, spacing);
    //let gridPositions = calculateGridPositions(grid, windowWidth, windowHeight, enemySize, spacing);//
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
    } */

    
	
})

go("startGame")
//go("game")








function getRandomInt(max) {
	return Math.floor(Math.random() * max);
  }