'use strict';

const WEEKDAY = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
const MONTH = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
const HOST = window.location.origin;

Number.prototype.pad = function(size) {
  var s = String(this);
  while (s.length < (size || 2)) {s = "0" + s;}
  return s;
};

var app = new Vue({
    el: '#app',
    data: {
      url: "",
      name: "",
      URL: null,
      ERROR: null
    },
    methods: {
      createURL: async function (url, name) {
        this.URL = null;
        this.ERROR = null;

        // Restart logo animation
        document.getElementById("loading-logo").contentDocument.documentElement.innerHTML += "";

        try {
          const response = await fetch(`${HOST}/${name}?url=${url}`, {method: 'POST', headers:{'Content-Type': 'application/json'}});
          if (!response.ok) {
            this.ERROR = {
              code: response.status,
              message: response.statusText.toUpperCase()
            };
            return;
          } 
          const body = await response.json();
          this.URL = body;
        } catch (error) {
          this.ERROR = {
            code: 666,
            message: error.message.toUpperCase()
          };
        }
      },
      updateURL: async function (url, name) {
        this.URL = null;
        this.ERROR = null;

        // Restart logo animation
        document.getElementById("loading-logo").contentDocument.documentElement.innerHTML += "";

        try {
          const response = await fetch(`${HOST}/${name}?url=${url}`, {method: 'PUT', headers:{'Content-Type': 'application/json'}});
          if (!response.ok) {
            this.ERROR = {
              code: response.status,
              message: response.statusText.toUpperCase()
            };
            return;
          } 
          const body = await response.json();
          this.URL = body;
        } catch (error) {
          this.ERROR = {
            code: 666,
            message: error.message.toUpperCase()
          };
        }
      },
      deleteURL: async function (name) {
        this.URL = null;
        this.ERROR = null;

        // Restart logo animation
        document.getElementById("loading-logo").contentDocument.documentElement.innerHTML += "";

        try {
          const response = await fetch(`${HOST}/${name}`, {method: 'DELETE', headers:{'Content-Type': 'application/json'}});
          if (!response.ok) {
            this.ERROR = {
              code: response.status,
              message: response.statusText.toUpperCase()
            };
            return;
          } 
          const body = await response.json();
          this.URL = body;
        } catch (error) {
          this.ERROR = {
            code: 666,
            message: error.message.toUpperCase()
          };
        }
      }
    },
    filters: {
      trim: function (value) {
        if (!value) return '';
        value = value.toString();
        return value.length >= 20 ? `${value.substring(0, 20)}...` : value;
      },
      formatDate: function (value) {
        if (!value) return '';
        value = new Date(Date.parse(value.toString()));
        return `${WEEKDAY[value.getDay()]}, ${value.getDate()} ${MONTH[value.getMonth()]} ${value.getFullYear()} ${value.getHours().pad(2)}:${value.getMinutes().pad(2)}`;
      }
    }
})
