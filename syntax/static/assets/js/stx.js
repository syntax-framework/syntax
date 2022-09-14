/**
 * Syntax Client
 *
 * Ver https://github.com/riot/dom-bindings
 */
(function (scriptTag) {
  // get server configuration
  let dataset = scriptTag.dataset;
  const serverEndpoint = dataset.endpoint;

  // @formatter:off
  const FILE            = "f"; // string
  const LINE            = "l"; // string
  const INITIALIZER     = "c"; // function($)
  const ELEMENTS        = "e"; // Array<key: elementIndex, value: string(#id|data-syntax-id)>
  const ATTRIBUTE_NAMES = "a"; // Array<key: attributeIndex, value: string>
  const EVENT_NAMES     = "n"; // Array<key: eventNameIndex, value: string>
  const EVENTS          = "o"; // Array<[elementIndex, eventNameIndex, expressionIndex]>
  const EXPRESSIONS     = "x"; // Array<key: expressionIndex, value: Function>
  const EXPORTS         = "z"; // Object, componente export methods

  // A writer applies the result of an expression to something (text, attribute, component, directive), has three forms
  //
  //  A) JS: Array<key: writerIndex, value: [elementIndex, expressionIndex]>
  //    Apply the result of an expression to an element ($(el).innerHtml = value)
  //
  //  B) JS: Array<key: writerIndex, value: [elementIndex, attributeIndex, expressionIndex]>
  //    Applies the result of the expression to an attribute ($(el).setAttribute(value))
  //
  //  C) JS: Array<key: writerIndex, value: [elementIndex, attributeIndex, [string, expressionIndex, string, ...]]>
  //    Apply the (dynamic) template to an attribute, allowing you to check for later changes to the attribute
  //    $(el).setAttribute(parse(template))
  const WRITERS         = "t"; // $writers

  // All watches. Represent expressions that will react when a variable changes.
  // JS: Array<key: _, value: [type, variableIndex, expressionIndex|writerIndex]>
  //    type 0 = action(expressionIndex)
  //    type 1 = schedule(writerIndex)
  const WATCHERS        = "w"; // Array<key: _, value: [type, variableIndex, expressionIndex|writerIndex]>



  // deprecado, remover da implementaçao
  // const ACTIONS         = "a"; // $actions            {Array<key: actionIndex, value: expressionIndex>} - uma ACTION é uma EXPRESSION
  // const EVENT_HANDLERS  = "h"; // $event_handlers    {Array<key: eventHandlerIndex, value: function>} é uma EXPRESSION

  // Instance attributes
  //const WATCHERS_VARS   = "v"; // $watchers_by_vars  {Array<key: variableIndex, value: watcherIndex[]>}
  // Rendered
  // const DYNAMICS = "d"
  // const STATIC = "s"
  // const COMPONENTS = "c"
  // const EVENTS = "e"
  // const REPLY = "r"
  // const TITLE = "t"
  // const TEMPLATES = "p"
  const PHASE = "P";
  const HANDLER = "H";
  const RENDER = "R";
  // @formatter:on

  // https://html.spec.whatwg.org/#boolean-attribute
  const BOOLEAN_ATTRIBUTES = [
    "allowfullscreen", "async", "autofocus", "autoplay", "checked", "controls", "default", "defer", "disabled",
    "formnovalidate", "ismap", "itemscope", "loop", "multiple", "muted", "nomodule", "novalidate", "open",
    "playsinline", "readonly", "required", "reversed", "selected", "truespeed",
  ];

  const STX = window.STX = {
    s: standalone, // standalone script
    b: bindToText,
    t: bindToAttr,
    channel: connectToChannel
  }

  /**
   * API com métodos utilitários disponíveis para uso pela instancia
   *
   * @type {{e: escape}}
   */
  const $_API = {
    e: escape
  }

  /**
   * Faz o escape de um conteúdo html
   *
   * @param value {string}
   * @return {string} o coteúdo html formatado
   */
  function escape(value) {
    return value
  }

  /* Watchers factories */

  function initChain() {
    // inicia um evento de atualização do estado.
    // mantém um stack da mudaná de estado para decidir se vai executar a próxima action ou não
    // chamadas recursivas devem ser armazenadas e analizadas para identificação de loops infinitos
    // em tempo de compilação o desenvolvedor recebe warnings sobre possíveis referencias cíclicas
    // a ordem de execução das ações é determinada pelo desenvolvedor, o compilador não utiliza DAG para otimizar
    // isso garante previsibilidade
    // sistema só faz scheduling para atualizar o DOM
  }

  function bindToExpression(elementIndex, expressionIndex) {
    return ($, ctx) => {
      let element = $(ctx[ELEMENTS][elementIndex]);
      let expression = ctx[EXPRESSIONS][expressionIndex];
      return {
        [PHASE]: RENDER,
        [HANDLER]: ($, ctx) => {
          element.h(expression())
        }
      }
    }
  }

  /**
   * Cria um watcher que altera o conteúdo de um elemento em tempo de renderização
   *
   * @param elementIndex
   * @param expressionIndex
   * @returns {(function(*, *))|*}
   */
  function bindToText(elementIndex, expressionIndex) {
    return ($, ctx) => {
      let element = $(ctx[ELEMENTS][elementIndex]);
      let expression = ctx[EXPRESSIONS][expressionIndex];
      return {
        [PHASE]: RENDER,
        [HANDLER]: ($, ctx) => {
          element.h(expression())
        }
      }
    }
  }

  /**
   * Cria um watcher que altera um atributo
   * @param elementIndex
   * @param attributeIndex
   * @param expressionIndex
   * @returns {function(*, *): {}}
   */
  function bindToAttr(elementIndex, attributeIndex, expressionIndex) {
    return ($, ctx) => {
      let element = $(ctx[ELEMENTS][elementIndex]);
      let expression = ctx[EXPRESSIONS][expressionIndex];
      return {
        [PHASE]: RENDER,
        [HANDLER]: ($, ctx) => {
          element.h(expression())
        }
      }
    }
  }

  function bindToAttrTpl(elementIndex, attributeIndex, expressionIndex) {
    return ($, ctx) => {

    }
  }

  // Create a watcher that bind element attribute to expression result
  // _$bind_prop(elementIndex, attributeIndex, expressionIdx )
  //   _$bind_prop(0, 0, 0) /* ${inputValue} */,
  //   _$bind_prop_tpl(0, 0, ['', 0, ''])

  /**
   * Registra o comportamento de um componente
   *
   * @param name
   * @param factory
   * @param dependencies
   */
  function component(name, factory, dependencies) {
    const $ = {
      p: onChangeParams
    }

    function onChangeParams(callback) {

    }
  }

  const standalones = []

  /**
   * Agenda a execução de um comportamento standalone, que não pertence a um componente específico.
   *
   * Uma script standalone é qualquer <script> existente dentro do html, usado para adicionar comportamentos a uma tag
   *
   * @param selector
   * @param factory
   */
  function standalone(selector, factory, dependencies) {
    // document.querySelector('#hQj7cqDhwKI')
    standalones.push({
      selector: selector,
      factory: factory,
      dependencies: dependencies
    });
  }


  /**
   * Faz a transformação do factory de um behavior em um construtor.
   *
   * Um factory de componente tem a mesma assinatura de uma standalone
   *
   * @param factory {(Object)=> Object}
   */
  function construct(factory) {

    const config = factory(STX)
    const file = config[FILE];
    const line = config[LINE];
    const elements = config[ELEMENTS]; // Array<key: elementIndex, value: string(#id|data-syntax-id)>
    const attributes = config[ATTRIBUTE_NAMES]; // Array<key: attributeIndex, value: string>
    const eventNames = config[EVENT_NAMES]; // Array<key: eventNameIndex, value: string>
    const events = config[EVENTS]; // Array<[elementIndex, eventNameIndex, expressionIndex]>

    const initializer = config[INITIALIZER]; // function($)

    // A writer applies the result of an expression to something (text, attribute, component, directive), has three forms
    //
    //  A) JS: Array<key: writerIndex, value: [elementIndex, expressionIndex]>
    //    Apply the result of an expression to an element ($(el).innerHtml = value)
    //
    //  B) JS: Array<key: writerIndex, value: [elementIndex, attributeIndex, expressionIndex]>
    //    Applies the result of the expression to an attribute ($(el).setAttribute(value))
    //
    //  C) JS: Array<key: writerIndex, value: [elementIndex, attributeIndex, [string, expressionIndex, string, ...]]>
    //    Apply the (dynamic) template to an attribute, allowing you to check for later changes to the attribute
    //    $(el).setAttribute(parse(template))
    const writers = config[WRITERS];


    // All watches. Represent expressions that will react when a variable changes.
    // JS: Array<key: _, value: [type, variableIndex, expressionIndex|writerIndex]>
    //    type 0 = action(expressionIndex)
    //    type 1 = schedule(writerIndex)
    // Array<key: _, value: [type, variableIndex, expressionIndex|writerIndex]>
    const watchers = config[WATCHERS];

    return Constructor

    /**
     * Construtor de uma instancia
     *
     * @param element {HTMLElement}
     * @constructor
     */
    function Constructor(element) {

      let cache_el = {}

      // Instance utility object
      let $ = function (selector, single) {
        if (Number.isInteger(selector)) {
          selector = elements[selector];
        }
        if (single) {
          let cached = cache_el[selector];
          if (!cached) {
            let el = $find(selector, element, single)
            if (el) {
              if (el.nodeName === "EMBED") {
                let text = document.createTextNode("");
                el.parentNode.insertBefore(text, el);
                el.parentNode.removeChild(el);
                el = text
              }
              cache_el[selector] = el
              cached = el;
            }
          }
          return cached
        }

        return $find(selector, element, single)
      }

      // copy common api
      Object.entries($_API).forEach(([key, value]) => {
        $[key] = value;
      });

      $.el = element;
      $.i = invalidate;
      $.c = onInputChange;


      /*
       * ORDEM DE INICIALIZAÇÃO
       *
       * 1. CONSTRUCTOR   Executa o construtor (`c: function($, STX){ }`) para obter Expressions, Lifecycle e API
       * 2. REFERENCES    Iniciado no momento da construção (@TODO: Mudar, deixar o controle para o framework)
       * 3. WATCHERS      Observa as mudanças nas variáveis
       * 4. EVENTS        Adiciona todos os eventos no DOM
       * 5. HOOKS         Dispara evento OnMount
       * *. WRITERS       Os Writers são criados conforme as variáveis sofrem alteração. Na primeira execução, todas as
       *                  variáveis são marcadas como alteradas para renderização
       */

      // @formatter:off
        let isUpdating = false;
        let isMounting = true;
        const dirty = new Set();                    // indices de variáveis que sofreram alteração
        const observers = [];                       // lista de funções que estão observando mudança em variáveis
        const writersScheduled = [];                // All writers that need to be executed
        const writers_instances = [];               // All writers that were initialized

        let onError =  function(trace, err) {
          console.error(new Error(`Whoops, something bad happened! file: ${file}, line: ${line}, trace: ${trace}, cause: ${err}`));
          console.error(err);
        }

        const instance = initializer($, STX);       // Initialize script
        const api = instance[EXPORTS];              // Instance API reference
        const expressions = instance[EXPRESSIONS];  // Array<key: expressionIndex, value: Function>

        // Life Cycle
        const onMount       = instance['a'] || noop;
        const beforeUpdate  = instance['b'] || noop;
        const afterUpdate   = instance['c'] || noop;
        const beforeRender  = instance['d'] || noop;
        const afterRender   = instance['e'] || noop;
        const onDestroy     = instance['f'] || noop;
        const onConnect     = instance['g'] || noop;
        const onDisconnect  = instance['h'] || noop;
        const onEvent       = instance['j'] || noop;
        onError             = instance['i'] || onError;
        // @formatter:on

      // Request Animation Frame
      let raf;

      if (watchers) {
        // All watches. Represent expressions that will react when a variable changes.
        // JS: Array<key: _, value: [type, variableIndex, expressionIndex|writerIndex]>
        //    type 0 = action(expressionIndex)
        //    type 1 = schedule(writerIndex)
        // Array<key: _, value: [type, variableIndex, expressionIndex|writerIndex]>
        watchers.forEach(([type, variableIndex, typeIndex]) => {
          let observer
          if (type === 0) {
            // type 0 = action(expressionIndex)
            const expression = expressions[typeIndex]
            observer = () => {
              expression();
            };
          } else {
            // type 1 = schedule(writerIndex)
            const writerIndex = typeIndex;
            observer = () => {
              // Agenda a escrita no dom
              if (!writersScheduled.includes(writerIndex)) {
                writersScheduled.push(writerIndex)
              }

              if (!raf) {
                raf = requestAnimationFrame(render)
              }
            };
          }

          if (!observers[variableIndex]) {
            observers[variableIndex] = [];
          }
          observers[variableIndex].push(observer);
        })
      }

      // Events
      if (events) {
        events.forEach(([elementIndex, eventNameIndex, expressionIndex]) => {
          let eventTarget = $(elements[elementIndex], true)
          if (eventTarget) {
            // @TODO: Options https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener
            eventTarget.addEventListener(eventNames[eventNameIndex], expressions[expressionIndex])
          }
        })
      }

      try {
        isMounting = false;
        onMount(); // hook
      } catch (e) {
        onError('onMount', e);
      }


      return api

      /**
       * Observa modificações nas variáveis e agenda atualizações no DOM
       *
       * @param variableIndex
       * @param o
       * @param n
       * @return {*}
       */
      function invalidate(variableIndex, o, n) {
        if (o !== o ? n === n : o !== n || ((o && typeof o === 'object') || typeof o === 'function')) {
          if (!isMounting && !isUpdating) {
            isUpdating = true;
            try {
              beforeUpdate(); // hook
            } catch (e) {
              onError('beforeUpdate', e);
            }

            queueMicrotask(() => {
              isUpdating = false;
              try {
                afterUpdate(); // hook
              } catch (e) {
                onError('afterUpdate', e);
              }
            });
          }
          dirty.add(variableIndex);
          dispatch(variableIndex);
        }
        return n;
      }

      /**
       * Notifica a mudança de uma variável
       *
       * @param variableIndex
       */
      function dispatch(variableIndex) {
        if (observers[variableIndex]) {
          observers[variableIndex].forEach(callback => {
            try {
              callback();
            } catch (e) {
              onError('dispatch', e);
            }
          });
        }
      }

      /**
       * Executa todas as escritas agendadas no DOM
       */
      function render() {
        try {
          beforeRender(); // hook
        } catch (e) {
          onError('beforeRender', e);
        }

        raf = undefined;
        writersScheduled.forEach((writerIndex) => {
          // get and exec writer
          try {
            getWriter(writerIndex)();
          } catch (e) {
            onError('render', e);
          }
        });
        writersScheduled.splice(0);

        try {
          afterRender(); // hook
        } catch (e) {
          onError('afterRender', e);
        }

        // @TODO: mudar estado de execução, verificar se existe tarefa pendente
      }

      /**
       * Cada instancia do componente faz a ligação do elemento com o writer
       *
       * @param writerIndex
       * @return {*}
       */
      function getWriter(writerIndex) {
        let writer = writers_instances[writerIndex]
        if (!writer) {
          // Faz a inicialização do writer
          // A writer applies the result of an expression to something (text, attribute, component, directive),
          // has three forms
          const [elementIndex, param2, param3] = writers[writerIndex];
          const element = $(elements[elementIndex], true);


          if (param3 === undefined) {
            //  A) JS: Array<key: writerIndex, value: [elementIndex, expressionIndex]>
            //    Apply the result of an expression to an element ($(el).innerHtml = value)
            const expression = expressions[param2];
            writer = () => {
              element.textContent = expression();
            }
          } else if (Number.isInteger(param3)) {
            //  B) JS: Array<key: writerIndex, value: [elementIndex, attributeIndex, expressionIndex]>
            //    Applies the result of the expression to an attribute ($(el).setAttribute(value))
            const attrName = attributes[param2];
            const expression = expressions[param3];

            if (attrName === 'value') {
              writer = () => {
                element.value = expression();
              }
            } else if (BOOLEAN_ATTRIBUTES.includes(attrName)) {
              writer = () => {
                if (!!expression()) {
                  element.setAttribute(attrName, attrName);
                } else {
                  element.removeAttribute(attrName)
                }
              }
            } else {
              writer = () => {
                element.setAttribute(attrName, expression());
              }
            }
          } else {
            //  C) JS: Array<key: writerIndex, value: [elementIndex, attributeIndex, [string, expressionIndex, string, ...]]>
            //    Apply the (dynamic) template to an attribute, allowing you to check for later changes to the attribute
            //    $(el).setAttribute(parse(template))
            // @TODO: Parse do template, verifica se modificou algo e atualiza o template
          }
          writers_instances[writerIndex] = writer;
        }

        return writer
      }

      /**
       * Two-way data-binding
       *
       * Specialized method for handling form elements (input, select, textarea).
       *
       * Automatically added when element does not have `onchange` and `oninput` events and is referencing a
       * scope variable (let, var) in the form `<input value="${myVariable}">`.
       *
       * Syntax adds the two events and applies the variable change.
       *
       * @param event {Event}
       * @param setterExpressionIndex {number}
       */
      function onInputChange(event, setterExpressionIndex) {
        let elem = event.target;
        let setValue = expressions[setterExpressionIndex];

        try {
          setValue(elem.value);
        } catch (e) {
          onError('onInputChange', e);
        }

        // @TODO: Form validation https://www.baeldung.com/javax-validation? Formatters?
        // @TODO: radio, checkbox, select-multiple, date, datetime-local, time, week, month, number, url, email, range
        // if (elem.type === 'radio' || elem.type === 'checkbox') { }
      }
    }

  }


  //-- DOM UTILS - START -----------------------------------------------------------------------------------------------

  /**
   * Faz a busca por um elemento no DOM
   *
   * @param selector {string}
   * @param scope {HTMLElement}
   * @param singleResult {boolean}
   * @return {NodeListOf<*>|*}
   */
  function $find(selector, scope, singleResult) {
    scope = scope || document
    let result = scope.querySelectorAll(selector)
    if (singleResult) {
      return result[0]
    }
    return Array.from(result)
  }

  //-- DOM UTILS - END -------------------------------------------------------------------------------------------------

  function getUUID() {
    return (
      String(1e7) + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, (c) => (
        Number(c) ^ (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (Number(c) / 4)))
      ).toString(16)
    );
  }


  function noop() {

  }

  /**
   * Copyright 2016 Andrey Sitnik <andrey@sitnik.ru>, https://github.com/ai/nanoevents/blob/main/LICENSE
   * @return {{emit(*, ...[*]): void, on(*, *): function(): void}|(function(): void)|*}
   */
  function createNanoEvents() {
    let bindingRef = 0;
    const events = {};
    return {
      emit(event, ...args) {
        let callbacks = events[event] || []
        for (let i = 0, length = callbacks.length; i < length; i++) {
          callbacks[i](...args)
        }
      },
      on(event, cb) {
        events[event]?.push(cb) || (events[event] = [cb])
        // off
        return () => {
          let callbacks = events[event]
          if (callbacks) {
            let idx = callbacks.indexOf(cb);
            if (idx >= 0) {
              callbacks.splice(idx, 1)
            }
          }
        }
      }
    }
  }

  //-- PUB SUB - START -------------------------------------------------------------------------------------------------


  // Single connection for entire application
  let connection

  function push(payload, timeout, maxRetries = 3) {
    let promise = new Promise((resolve, reject) => {
      const fetchWithRetries = async (retries) => {
        try {
          return await fetch(serverEndpoint, {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            credentials: 'same-origin',
            headers: {
              'Accept': 'application/json',
              'Content-Type': 'application/json',
            },
            redirect: 'follow',
            referrerPolicy: 'no-referrer',
            body: JSON.stringify(payload)
          });
        } catch (error) {
          if (retries < maxRetries) {
            return fetchWithRetries(retries + 1);
          }
          // max retries exceeded
          throw error;
        }
      }
      return fetchWithRetries(0)
    })

    // if the timeout option is specified, race the fetch call
    if (timeout) {
      promise = Promise.race([
        promise,
        new Promise((_, reject) => {
          // throw on timeout
          setTimeout(() => reject("TIMEOUT"), timeout)
        })
      ]);
    }

    return promise;
  }

  class Channel {

    constructor(topic, params, sse) {
      this.seq = 0;
      this.topic = topic
      this.timeout = 5000; // @TODO: from server config
      this.events = createNanoEvents()

      let listener = (e) => {
        console.log(this, e.data)
      };
      sse.addEventListener(topic, listener)
      this.onClose(() => {
        sse.removeEventListener(listener)
      });
    }

    onClose(callback) {
      this.on("stx_close", callback);
    }

    on(event, callback) {
      return this.events.on(event, callback);
    }

    /**
     * Sends a message `event` to syntax with the payload `payload`.
     *
     * Syntax receives this in the `handle_in(event, payload, socket)` function. if syntax replies or it times out (default 10000ms),
     * then optionally the reply can be received.
     *
     * @example
     * channel.push("event")
     *   .then(() => console.log("Submitted"))
     *   .catch(err => console.log("Syntax errored", err))
     *
     * @param {string} event
     * @param {Object} payload
     * @param {number} timeout
     * @param {number} maxRetries
     * @returns {Promise<Response>}
     */
    push(event, payload, timeout = this.timeout, maxRetries) {
      return push({t: this.topic, e: event, s: this.seq++, p: payload || {}}, timeout, maxRetries)
    }

    /**
     * Leaves the channel
     *
     * Unsubscribes from server events, and instructs channel to terminate on server
     * @param timeout
     */
    close(timeout = this.timeout) {
      // return this.events.on(event, callback);
    }
  }

  /**
   *
   * @param topic
   * @param params
   */
  function connectToChannel(topic, params) {
    if (!connection) {
      connection = (() => {

        let channels = [];

        // https://developer.mozilla.org/en-US/docs/Web/API/EventSource
        let sse = new EventSource(serverEndpoint);

        sse.onmessage = (event) => {
          let msg = JSON.parse(event.data);
          console.log("New Live Event", msg)
        }

        sse.onerror = (event) => {
          console.log('onerror', event)
        }

        sse.onopen = (event) => {
          console.log('onopen', event)
        }

        /**
         *
         * @param topic
         * @param params
         */
        function channel(topic, params) {
          // channel:topic
          topic = topic.replaceAll(/[^A-Za-z0-9\-_:]/g, '').replaceAll(/[:]+/g, ':')
          if (topic === '') {
            throw new Error('Invalid topic name')
          }
          let parts = topic.split(':')

          if (parts.length > 2) {
            throw new Error('Invalid topic name: ' + topic)
          }

          let channel = new Channel(topic, params, sse)

          channels.push(channel)
          channel.onClose(() => {
            let idx = channels.indexOf(channel);
            if (idx >= 0) {
              channels.splice(idx, 1)
            }
          })

          return channel
        }

        return {
          channel: channel
        }
      })()
    }

    return connection.channel(topic, params)
  }

  //-- PUB SUB - END ---------------------------------------------------------------------------------------------------

  /**
   * Starts the Syntax Client
   */
  function initialize() {
    // initialize components
    standalones.forEach((config) => {
      let element = $find(config.selector, null, true)
      if (element) {
        const constructor = construct(config.factory)
        constructor(element)
      }
    })

    let channel = STX.channel("xpto:beribecanta", {nome: 'albumina'});
    console.log(channel)
  }

  if (document.readyState === 'loading') {
    document.addEventListener("DOMContentLoaded", initialize);
  } else {
    initialize();
  }

})(document.currentScript)

