const turn_box = document.getElementById('turn_order');
var turn_order = [{name: 'top of the round', initiative: -99}];
const counter = document.getElementById('turn_counter');
var turn_counter = 1;



function roll() {
    let stat = document.getElementById('prof_list').value;
    let characters_opt = document.getElementById('char_list').selectedOptions;
    let characters = [];
    
    for (i = 0; i< characters_opt.length; i++) {
        characters.push(characters_opt[i]['value']);
    }
    
    var output_box = document.querySelector('#output_box');
    var anchor = document.querySelector('#anchor');

    function roll_i(character, stat){
        let b = Math.floor(Math.random() * 20) + 1;
        let stats = data;
        let m = stats[character][stat] + b;
        let training = stats[character][`${stat}Train`];
        if(training !== null){    
            if(b === 20){
                var result = `${character} ${stat}: ${m} (${training}) Nat 20!!!`;
            } else if(b === 1){
                var result = `${character} ${stat}: ${m} (${training}) Nat 1!!!`;
            } else {
                var result = `${character} ${stat}: ${m} (${training}) `;
            }

            let output =  document.createElement('div');
            output.className = 'text_output';
            output.innerText = result;
            output_box.insertBefore(output, anchor);
        }
    }
    
    characters.forEach(d => roll_i(d, stat));
}



function add_turn(){
    let turn_name = document.querySelector('#turn_name').value;
    let turn_num = document.querySelector('#turn_num').value;

    turn_num = Number(turn_num);

    if(turn_name !== '' && turn_num !== '' && isNaN(turn_num) === false){
        let turn = {name: turn_name, initiative: turn_num};

        turn_order.push(turn);
    }

    document.getElementById('turn_name').value = '';
    document.getElementById('turn_num').value = '';

    if(document.contains(document.getElementById('turn_list'))){
        document.getElementById('turn_list').remove();
    }
    
    var turn_list = document.createElement('div');
    turn_list.id = 'turn_list';
    turn_box.appendChild(turn_list);

    turn_order.forEach(i => post_turn(i));
}



function sort(){
    turn_order.sort(function(a, b){return b.initiative - a.initiative});

    if(document.contains(document.getElementById('turn_list'))){
        document.getElementById('turn_list').remove();
    }
    
    var turn_list = document.createElement('div');
    turn_list.id = 'turn_list';
    turn_box.appendChild(turn_list);

    turn_order.forEach(i => post_turn(i));
}



function clear_initiative(){
    turn_order = [{name: 'top of the round', initiative: -99}];

    if(document.contains(document.getElementById('turn_list'))){
        document.getElementById('turn_list').remove();
    }

    if(document.contains(document.getElementById('count'))){
        document.getElementById('count').remove();
    }
    turn_counter = 1;
    let count = document.createElement('div');
    count.innerText = turn_counter;
    count.id = 'count';
    counter.appendChild(count);
}



function next_turn(){

    if(turn_order.length > 1){
        if(turn_order[1]['name'] === 'top of the round'){
            if(document.contains(document.getElementById('count'))){
                document.getElementById('count').remove();
            }
            
            turn_counter +=1;
            let count = document.createElement('div')
            count.innerText = turn_counter;
            count.id = 'count';
            turn_order.push(turn_order.shift());

            counter.appendChild(count);
        }
    }

    turn_order.push(turn_order.shift());

    if(document.contains(document.getElementById('turn_list'))){
        document.getElementById('turn_list').remove();
    }
    
    var turn_list = document.createElement('div');
    turn_list.id = 'turn_list';
    turn_box.appendChild(turn_list);

    turn_order.forEach(i => post_turn(i));
}

function post_turn(i){
    if(i.name !== 'top of the round'){    
        let text = `${i.initiative}  ${i.name}`;
        let line = document.createElement('div');
        line.className = 'turn';
        line.innerText = text;

        turn_list.appendChild(line);
    }
}