import React from 'react';
import useState from 'react'
import logo from './logo.svg';
import './App.css';
import {Button} from "rbx";

class App  extends React.Component{

  // {} means its an object
  // {} = Objects
  // [] = Arrays
  //useState = {output: {}}

  constructor(props){
    super(props);
    this.state = {
      data: null,
      play1: null,
      play2: null,
      testNumber: 0
    };
  }


  componentDidMount(){
    //const [output, setOuput] = useState({});

    // basically calling a get onto the url there
    fetch('http://localhost:8080')
          //.then(res => this.setState({output: res.json}))
          .then((res) => {
            //console.log(JSON.parse(res.body))
            return res.json()
          })
          .then((data) => {
            //console.log(data);
            //console.log(JSON.parse('{"play1":{"name":"Noah","Hand":[{"Action":"MILL","Count":3},{"Action":"DRAW","Count":1},{"Action":"MILL","Count":3}],"Drawing":1,"Milling":0},"play2":{"name":"roBob","Hand":[{"Action":"SKIP","Count":2},{"Action":"DRAW","Count":1},{"Action":"DRAW","Count":1}],"Drawing":1,"Milling":0},"deck":[{"Action":"DRAW","Count":1},{"Action":"DRAW","Count":1},{"Action":"SKIP","Count":5},{"Action":"DRAW","Count":1},{"Action":"SKIP","Count":5},{"Action":"DRAW","Count":1},{"Action":"DRAW","Count":10},{"Action":"DRAW","Count":1},{"Action":"DRAW","Count":1},{"Action":"DRAW","Count":1},{"Action":"MILL","Count":3},{"Action":"MILL","Count":3},{"Action":"SKIP","Count":8},{"Action":"MILL","Count":6},{"Action":"MILL","Count":3},{"Action":"SKIP","Count":2},{"Action":"SKIP","Count":2},{"Action":"SKIP","Count":2},{"Action":"MILL","Count":6},{"Action":"SKIP","Count":2},{"Action":"DRAW","Count":5},{"Action":"DRAW","Count":5},{"Action":"DRAW","Count":5},{"Action":"SKIP","Count":5}],"TESTING":"Is this thing on?"}'));
            this.setState({
              data: data
            });
            console.log("3",this.state.data);
          })
          /*.then((testVariable) => {
            this.setState({ contacts: testVariable })
          })*/
          .catch(console.log)
  }

  render (){
    const { mydata } = this.state;

    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <p>
            Edit <code>src/App.js</code> and save to reload.
          </p>
          <Button>Hi</Button>
          {this.state.data && <p>Look at this = {this.state.data.play1.name}</p>}
          <a
            className="App-link"
            href="https://reactjs.org"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn React
          </a>
        </header>
      </div>
    );
  }

}

export default App;
