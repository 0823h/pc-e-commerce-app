import React, { useEffect, useState } from "react";

import axios from "axios";

const ProductList = (props) => {
  // eslint-disable-next-line
  const [productList, setProductList] = useState(null);
  axios
    .get(`localhost:8080/api/products`)
    .then((result) => {
      this.setProductList(result.data);
      console.log("abclog:" + result.data);
    })
    .catch((error) => console.log(error));
  console.log("hello world");
  return <div>Hello: {productList} HelloHello</div>;
};

export default ProductList;
