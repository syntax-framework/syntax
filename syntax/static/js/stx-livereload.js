(function (scriptTag) {
  // get server configuration
  let dataset = scriptTag.dataset;
  const endpoint = dataset.endpoint;
  const interval = Number.parseInt(dataset.interval);
  const reloadPageOnCss = dataset.reloadPageOnCss === 'true';

  function buildFreshUrl(link) {
    const date = Math.round(Date.now() / 1000).toString();
    const url = link.href.replace(/(\&|\\?)vsn=\d*/, '');
    const newLink = document.createElement('link');
    const onComplete = function () {
      if (link.parentNode !== null) {
        link.parentNode.removeChild(link);
      }
    };

    newLink.onerror = onComplete;
    newLink.onload = onComplete;
    link.setAttribute('data-pending-removal', '');
    newLink.setAttribute('rel', 'stylesheet');
    newLink.setAttribute('type', 'text/css');
    newLink.setAttribute('href', url + (url.indexOf('?') >= 0 ? '&' : '?') + 'vsn=' + date);
    link.parentNode.insertBefore(newLink, link.nextSibling);

    return newLink;
  }

  function repaint() {
    if (navigator.userAgent.toLowerCase().indexOf('chrome') > -1) {
      setTimeout(function () {
        document.body.offsetHeight;
      }, 25);
    }
  }

  function cssStrategy() {
    let selectors = 'link[rel=stylesheet]:not([data-no-reload]):not([data-pending-removal])';
    [].slice
      .call(window.parent.document.querySelectorAll(selectors))
      .filter(function (link) {
        return link.href
      })
      .forEach(function (link) {
        buildFreshUrl(link)
      });

    repaint();
  }

  function pageStrategy() {
    window.location.reload();
  }

  const reloadStrategies = {
    css: reloadPageOnCss ? pageStrategy : cssStrategy,
    page: pageStrategy
  };

  // Use SSE
  //
  // https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events
  new EventSource(endpoint).onmessage = (event) => {
    let msg = JSON.parse(event.data);
    console.log("New Live Reload Event", msg)
    setTimeout(function () {
      // const reloadStrategy = reloadStrategies[msg.asset_type] || reloadStrategies.page;
      // reloadStrategy();
    }, interval);
  }

  /*
  Part of the code obtained from https://github.com/phoenixframework/phoenix_live_reload

  # MIT License

  Copyright (c) 2014 Chris McCord

  Permission is hereby granted, free of charge, to any person obtaining
  a copy of this software and associated documentation files (the
  "Software"), to deal in the Software without restriction, including
  without limitation the rights to use, copy, modify, merge, publish,
  distribute, sublicense, and/or sell copies of the Software, and to
  permit persons to whom the Software is furnished to do so, subject to
  the following conditions:

  The above copyright notice and this permission notice shall be
  included in all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
  MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
  NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
  LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
  OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
  WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
})(document.currentScript)
