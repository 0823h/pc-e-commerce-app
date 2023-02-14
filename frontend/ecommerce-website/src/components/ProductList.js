// import React from 'react';

// import axios from 'axios';

// export default class ProductList extends React.Component {
//   state = {
//     persons: []
//   }

//   componentDidMount() {
//     axios.get(`http://localhost:8080/api/products`, {
//       headers: {"Access-Control-Allow-Origin": "*"},
//       responseType: "json"
//     })
//       .then(res => {
//         const persons = res.data;
//         this.setState({ persons });
//       })
//       .catch(error => console.log(error));
//   }

//   render() {
//     return (
//       <ul>
//         { this.state.persons.map(person => <li>{person.name}</li>)}
//       </ul>
//     )
//   }
// }

import React from 'react';

import axios from 'axios';

export default class ProductList extends React.Component {
  state = {
    persons: []
  }

  async componentDidMount() {
    await axios.get(`http://localhost:8080/api/products`)
      .then(res => {
        const persons = res.data;
        this.setState({ persons });
      })
      .catch(error => console.log(error));
  }

  render() {
    return (
      <ul>
        { this.state.persons.map(person => <li>{person.name}</li>)}
      </ul>
    )
  }
}