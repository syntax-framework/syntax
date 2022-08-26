/**
 * Syntax Client
 *
 * Ver https://github.com/riot/dom-bindings
 */
(function () {
  // @formatter:off
  const FILE            = "f"; // $file              {string}
  const LINE            = "l"; // $line              {string}
  const INITIALIZER     = "i"; // initializer        {function($)}
  const ELEMENTS        = "e"; // $elements          {string[]}
  const ATTRIBUTE_NAMES = "a"; // $attribute_names   {string[]}                Ex. `{ t: ['value', 'href', 'src'] }`
  const EVENT_NAMES     = "n"; // $event_names       {string[]}                Ex. `{ o: ['submit', 'change'] }`
  const EVENTS          = "o"; // $events            {Array<[elementIndex, eventNameIndex, expressionIndex]>}

  // todas as expressões existentes
  // eventos, interpolaçoes, dual bindings (gerados quando <elemento value="${variable}">)
  const EXPRESSIONS     = "x"; // $expressions        {Array<key: expressionIndex, value: Function>}

  // Um writer aplica o resultado de uma expression em algo (texto, atributo, componente, directiva), tem dois formatos
  //
  //  A) {Array<key: writerIndex, value: [elementIndex, expressionIndex]>}
  //      Aplica o resultado de uma expressão em um elemento (innerHtml)
  //
  //  B) {Array<key: writerIndex, value: [elementIndex, attributeIndex, expressionIndex]>}
  //      Aplica o resultado da expressao em um atributo
  //
  //  C) {Array<key: writerIndex, value: [elementIndex, attributeIndex, [static, expressionIndex, static, ...]]>}
  //      Aplica o template (dinamico) em um atributo, permitindo verificar mudanças posteriores no atributo
  const WRITERS         = "t"; // $writers

  // Observa alteração de variáveis, disparando açoes ou agendando scrita
  //  type: 0 = action, 1 = schedule
  const WATCHERS        = "w"; // $watchers          {Array<key: _, value: [type, variableIndex, expressionIndex|writerIndex]>}



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

  const STX = {
    s: standaloneScript, // standalone script
    b: bindToText,
    t: bindToAttr,
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

  function component() {
    const $ = {
      p: onChangeParams
    }

    function onChangeParams(callback) {

    }
  }

  /**
   * Executa um script standalone, que não pertence a um componente específico
   *
   * @param element
   * @param factory
   */
  function standaloneScript(element, factory) {

  }

  window.STX = STX;
})()

