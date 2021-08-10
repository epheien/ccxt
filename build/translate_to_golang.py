#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# pip3 install esprima pyhumps

import os
import sys
import traceback
import re
import humps
import subprocess
import shlex

import esprima
import json


EX_NAME = ''
CODE_INFO = {}
PWD = os.path.dirname(__file__)


def get_func_code_map(code_str):
    return []


def init_code():
    return '''
package {ex}

import (
	. "ccxt-master/go/base"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type {ex.capitalize()} struct {
	Exchange
}

    '''


JS_FOLD = 'js'

EX_EXTRA_FUNC = {
    'kucoin': ['fetchOrdersByStatus']
}
FUNC_LIST = [
    'sign', 'fetchOrderBook', 'fetchOpenOrders', 'cancelOrder',
    'createOrder', 'fetchOrder', 'parseOrder', 'fetchBalance',
    'fetchOrdersByStatus', 'fetchOrdersByState', 'fetchMarkets',
    'fetchCurrencies', 'handleErrors', 'fetchAccounts',
    'parseOrderStatus', 'fetchOrdersByStates',
    'ParseOrderSide', 'parseMarkets', 'parseMarket',
    'fetchMarketsByType', 'getPathAuthenticationType', 'parseOrderSide',
    'parseBalanceByType', 'parseAccountBalance', 'parseMarginBalance', 'parseFuturesBalance', 'parseSwapBalance'
]
JS_PATCH_FOR_GOLAGNG_TRANSLATE = {
    'kucoin': {
        'sign': {
            "let amount = this.safeFloat (order, 'size');": "amount = xdfsa",
        }
    }
}

def read_describe(str_code):
    x = re.findall(r"super.describe \(\), ({.*?})\);", str_code, re.MULTILINE|re.DOTALL)[0]

    xx = ''
    for line in x.split('\n'):
        if '// ' in line:
            line = line[:line.find('//')]
        if 'this.' in line:
            continue
        xx += line + '\n'

    x = xx
    x = re.sub(r': ([A-Z].+),', r': "\1",', x)
    x = x.replace('undefined', 'None')
    x = x.replace('true', 'True')
    x = x.replace('false', 'False')
    a = eval(x)
    return a


def format_js_code_for_espima_analysis(func, block):
    block = block.replace('await ', '')
    block = block.replace('async ', 'function ')
    if 'function' not in block.split('\n')[0]:
        block = 'function' + block
    return block


def format_describe_func(desc):
    return f'''
    func (self *{EX_NAME}) Describe() []byte {{
	return []byte(`{json.dumps(desc, indent=4)}`)
	}}
	'''

def split_func_list(str_code):
    res = str_code.split('\n\n')
    if res:
        tmp = res[-1]
        res[-1] = tmp.replace('};', '')
    return res


def read_func(str_code):
    result = {}

    for block in split_func_list(str_code):
        for func in FUNC_LIST:
            if f'{func} ' in block.split('\n')[0]:
                result[func] = format_js_code_for_espima_analysis(func, block)
                # print(format_js_code_for_espima_analysis(func, block))

    return result


def read_code_str(ex):
    p = os.path.join('..', 'js', f'{ex}.js')
    with open(p, encoding='UTF-8') as f:
        str_code = f.read()

    return {
        'describe': read_describe(str_code),
        'func': read_func(str_code),
    }


FUNC_ARG_MAP = {
    'createOrder': 'symbol string, typ string, side string, amount float64, price float64, params map[string]interface{}',
    'fetchBalance': 'params map[string]interface{}',
    'cancelOrder': 'id string, symbol string, params map[string]interface{}',
    'fetchOpenOrders': 'symbol string, since int64, limit int64, params map[string]interface{}',
    'fetchOrdersByStatus': 'status string, symbol string, since int64, limit int64, params map[string]interface{}',
    'fetchOrdersByState': 'status string, symbol string, since int64, limit int64, params map[string]interface{}',
    'fetchOrdersByStates': 'states string, symbol string, since int64, limit int64, params map[string]interface{}',
    'fetchOrderBook': 'symbol string, limit int64, params map[string]interface{}',
    'fetchOrder': 'id string, symbol string, params map[string]interface{}',
    'sign': 'path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}',
    'parseOrder': 'order interface{}, market interface{}',
    'parseOrders': 'status string, symbol string, since int64, limit int64, params map[string]interface{}',
    'fetchMarkets': 'params map[string]interface{}',
    'fetchCurrencies': 'params map[string]interface{}',
    'handleErrors': 'httpCode int64, reason string, url string, method string, headers interface{}, body string, response interface{}, requestHeaders interface{}, requestBody interface{}',
    'fetchAccounts': 'params map[string]interface{}',
    'parseOrderStatus': 'status string',
    'ParseOrderSide': 'side string',
    'parseMarket': 'market interface{}',
    'parseMarkets': 'markets []interface{}',
    'fetchMarketsByType': 'typ string, params map[string]interface{}',
    'getPathAuthenticationType': 'path string',
    'parseOrderSide': 'side string',
    'parseBalanceByType': 'typ string, response interface{}',
    'parseAccountBalance': 'response interface{}',
    'parseMarginBalance': 'response interface{}',
    'parseFuturesBalance': 'response interface{}',
    'parseSwapBalance': 'response interface{}',
}
RETURN_MAP = {
    'createOrder': 'result *Order, err error',
    'fetchBalance': 'balanceResult *Account, err error',
    'cancelOrder': 'response interface{}, err error',
    'fetchOrdersByStatus': 'orders interface{}',
    'fetchOrdersByState': 'orders interface{}',
    'fetchOrdersByStates': 'orders interface{}',
    'fetchOpenOrders': 'result []*Order, err error',
    'fetchOrderBook': 'orderBook *OrderBook, err error',
    'fetchOrder': 'result *Order, err error',
    'sign': 'ret interface{}',
    'parseOrder': 'result map[string]interface{}',
    'fetchMarkets': '[] interface{}',
    'fetchCurrencies': 'map[string]interface{}',
    'handleErrors': '',
    'fetchAccounts': '[]interface{}',
    'parseOrderStatus': 'string',
    'ParseOrderSide': 'string',
    'parseMarket': 'interface{}',
    'parseMarkets': '[] interface{}',
    'fetchMarketsByType': '[] interface{}',
    'getPathAuthenticationType': 'string',
    'parseOrderSide': 'string',
    'parseBalanceByType': 'interface{}',
    'parseAccountBalance': 'map[string]interface{}',
    'parseMarginBalance': 'map[string]interface{}',
    'parseFuturesBalance': 'map[string]interface{}',
    'parseSwapBalance': 'map[string]interface{}',
}
PANIC_DEAL_FUNC = [o.lower() for o in [
    'fetchOrderBook', 'fetchOpenOrders', 'cancelOrder',
    'createOrder', 'fetchOrder', 'fetchBalance',
]]
NIL_MAP = {
    'string': '""',
    'int64': '0',
    'error': 'nil'
}
ERROR_RETURN_FUNCS = [o.lower() for o in []]
SIDE = None


def get_return(func_name):
    return RETURN_MAP.get(func_name, None)


def get_arg(func_name):
    default = '(interface{})'
    return FUNC_ARG_MAP.get(func_name, default)


DEFAULT_FUNC_ARGS = {
    'SafeStringLower'.lower(): (3, '""'),
    'SafeString2'.lower(): (4, '""'),
    'SafeInteger2'.lower(): (4, 0),
    'SafeFloat2'.lower(): (4, 0.),
    'SafeValue'.lower(): (3, 'nil'),
    'SafeString'.lower(): (3, '""'),
    'ApiFunc'.lower(): (4, 'nil'),
    'SafeFloat'.lower(): (3, 0),
    'SafeInteger'.lower(): (3, 0),
}


def ThisExpression(syntax, info={}):
    return 'self'


def MemberExpression(syntax, info={}):
    info['error_check'] = False
    method_name = call_func_by_syntax(syntax.property)
    if syntax.object.type == 'ThisExpression' and syntax.property.type == 'Identifier' and method_need_check_err(method_name):
        info['error_check'] = True

    if syntax.object.type == 'ThisExpression':
        method_name = humps.pascalize(method_name)

    obj = call_func_by_syntax(syntax.object)
    if syntax.object.type == 'ThisExpression':
        return f'{obj}.{method_name}'
    else:
        m = {
            'toString': f'fmt.Sprintf("%v", {obj})',
            'length': f'self.Length({obj})',
        }
        default = f'{obj}[{method_name}]'

        if syntax.property.type == 'Identifier' and syntax.property.name == 'push':
            if 'arg_str' in info:
                return f'{obj} = append({obj}, {info["arg_str"]})'
        if syntax.property.type == 'Identifier' and syntax.property.name == 'split':
            if 'arg_str' in info:
                return f'strings.Split({obj}, {info["arg_str"]})'

        # if info.get('pre') in ['right', 'init', 'test']:
        if info.get('pre') != 'left':
            default = f'self.Member({obj}, {method_name})'
        if info.get('pre') in ['left']:
            default = f'self.SetValue({obj}, {method_name})'

        return m.get(syntax.property.name, default)


def CallExpression(syntax, info={}):
    arg_str = ','.join([call_func_by_syntax(arg, info) for arg in syntax.arguments])

    # default arg
    if syntax.callee.property and syntax.callee.property.name.lower() in DEFAULT_FUNC_ARGS:
        arg_info = DEFAULT_FUNC_ARGS[syntax.callee.property.name.lower()]
        arg_str += f', {arg_info[1]}' * (arg_info[0] - len(syntax.arguments))
    info.update({'arg_str': arg_str})

    pre_part = f'{call_func_by_syntax(syntax.callee, info)}'

    if syntax.callee.type == 'MemberExpression' and syntax.callee.object.type == 'Identifier' and syntax.callee.object.name == 'Object':
        if syntax.callee.property.name == 'keys':
            return f'reflect.ValueOf({call_func_by_syntax(syntax.arguments[0])}).MapKeys()'

    # toString -> fmt.Sprintf
    if '(' in pre_part and ')' in pre_part:
        return pre_part

    call_str = pre_part
    api_func_keys = [capital_first(o) for o in CODE_INFO['describe']['api'].keys()]
    api_func_pattern = f"self\.({'|'.join(api_func_keys)})"
    if re.findall(api_func_pattern, call_str):
        info1 = DEFAULT_FUNC_ARGS['apifunc']
        arg_str += f', {info1[1]}' * (info1[0] - len(syntax.arguments) - 1)
        return f'self.ApiFunc("{syntax.callee.property.name}", {arg_str})'
    # xx.split('_') -> string.Split(xx, '_')
    elif not call_str.endswith(')'):
        return f'{call_str}({arg_str})'
    else:
        return call_str


def ExpressionStatement(syntax, info={}):
    info['error_check'] = False
    s = call_func_by_syntax(syntax.expression, info)

    if info.get('error_check'):
        return f'_, err = {s}\n if err != nil {{\n return \n}}'
    else:
        return s


def method_need_check_err(name, info={}):
    if any(name.lower().startswith(o) for o in ERROR_RETURN_FUNCS):
        return True
    return False


def VariableDeclarator(syntax, info={}):
    left = call_func_by_syntax(syntax.id)
    operator = ':='
    if left in ARG_TYPE or left in RET_TYPE:
        operator = '='

    info = {'pre': 'init'}
    if syntax.init.type == 'Identifier' and syntax.init.name == 'undefined':
        return f'var {left} interface{{}}'
    if syntax.id.type == 'ArrayPattern':
        return f'{left} {operator} self.Unpack{len(syntax.id.elements)}({call_func_by_syntax(syntax.init)})'

    right = call_func_by_syntax(syntax.init, info)

    if info.get('error_check'):
        return f'{left}, err {operator} {right}\nif err != nil {{\nreturn nil, err\n}}'
    else:
        return f'{left} {operator} {right}'


def VariableDeclaration(syntax, info={}):
    return '\n'.join([call_func_by_syntax(o) for o in syntax.declarations])


def ObjectExpression(syntax, info={}):
    k_v = ''
    for one in syntax.properties:
        k_v += f'\n{call_func_by_syntax(one.key)}: {call_func_by_syntax(one.value)},'
    ret = f'map[string]interface{{}}{{{k_v}\n}}'
    return ret


def Literal(syntax, info={}):
    if "'" in syntax.raw:
        val = syntax.raw.replace('"', '')
        val = val.replace("'", '')
        val = f'"{val}"'
    else:
        val = f'{syntax.value}'
    return val


def Identifier(syntax, info={}):
    m = {
        'undefined': 'nil',
        'type': 'typ',
    }
    return m.get(syntax.name, syntax.name)


def BinaryExpression(syntax, info={}):
    operator = syntax.operator

    if syntax.right.name == 'undefined':
        if operator == '===':
            return f'self.TestNil({call_func_by_syntax(syntax.left)})'
        if operator == '!==':
            return f'!self.TestNil({call_func_by_syntax(syntax.left)})'

    if operator == '===':
        operator = '=='
    if operator == '!==':
        operator = '!='
    if operator == 'in':
        # return f'_, ok := {call_func_by_syntax(syntax.right)}[{call_func_by_syntax(syntax.left)}]; ok'
        return f'self.InMap({call_func_by_syntax(syntax.left)}, {call_func_by_syntax(syntax.right)})'
    left = call_func_by_syntax(syntax.left)
    right = call_func_by_syntax(syntax.right)

    if left in ARG_TYPE and right == 'nil':
        right = NIL_MAP.get(ARG_TYPE[left], 'nil')

    return f'{left} {operator} {right}'


def call_func_by_syntax(syntax, info={}):
    try:
        return globals()[syntax.type](syntax, info)
    except Exception as e:
        print(traceback.format_exc())
        print(syntax)
        raise


def AssignmentExpression(syntax, info={}):
    operator = syntax.operator
    if syntax.left.type == 'MemberExpression':
        arg1 = call_func_by_syntax(syntax.left.object)
        arg2 = call_func_by_syntax(syntax.left.property)
        arg3 = call_func_by_syntax(syntax.right, {"pre": "right"})
        return f'self.SetValue({arg1}, {arg2}, {arg3})'
    else:
        left = call_func_by_syntax(syntax.left, {"pre": "left"})
        right = call_func_by_syntax(syntax.right, {"pre": "right"})
        return f'{left} {operator} {right}'


def BlockStatement(syntax, info={}):
    lines = '\n'.join([call_func_by_syntax(block) for block in syntax.body])
    return f'{{\n{lines}}}'


def IfStatement(syntax, info={}):
    ret = f'if {call_func_by_syntax(syntax.test)} '

    ret = f'if self.ToBool({call_func_by_syntax(syntax.test)}) '
    ret += f'{call_func_by_syntax(syntax.consequent)}'

    if syntax.alternate:
        ret += (
            f' else '
            f'{call_func_by_syntax(syntax.alternate)}'
        )

    return ret


def ConditionalExpression(syntax, info={}):
    return f'self.IfThenElse(self.ToBool({call_func_by_syntax(syntax.test, {"pre": "test"})}), {call_func_by_syntax(syntax.consequent)}, {call_func_by_syntax(syntax.alternate)})'


def ArrayPattern(syntax, info={}):
    return ', '.join([call_func_by_syntax(e) for e in syntax.elements])


def UnaryExpression(syntax, info={}):
    return f'{syntax.operator}self.ToBool({call_func_by_syntax(syntax.argument)})'


def ReturnStatement(syntax, info={}):
    if not syntax.argument:
        return f'return'

    if info.get('return') and ',' in info['return']:
        return f'return {call_func_by_syntax(syntax.argument)}, nil'
    else:
        return f'return {call_func_by_syntax(syntax.argument)}'


def UpdateExpression(syntax, info={}):
    return f'{call_func_by_syntax(syntax.argument)}{syntax.operator}'


def ForStatement(syntax, info={}):
    return f'for {call_func_by_syntax(syntax.init)}; {call_func_by_syntax(syntax.test)}; {call_func_by_syntax(syntax.update)} {call_func_by_syntax(syntax.body)}'


def ThrowStatement(syntax, info={}):
    return f'self.RaiseException("{syntax.argument.callee.name}", {call_func_by_syntax(syntax.argument.arguments[0])})'


def LogicalExpression(syntax, info={}):
    return f'{call_func_by_syntax(syntax.left)} {syntax.operator} {call_func_by_syntax(syntax.right)}'


def ArrayExpression(syntax, info={}):
    elements = ','.join([call_func_by_syntax(o) for o in syntax.elements])
    return f'[]interface{{}}{{{elements}}}'


def capital_first(s, info={}):
    return s[0].upper() + s[1:]


FUNC_LINES = 99
ARG_TYPE = dict()
RET_TYPE = dict()
def FunctionDeclaration(syntax, info={}):
    result_code = []
    func_name = syntax.id.name

    global ARG_TYPE, RET_TYPE

    func_arg_str = get_arg(func_name)
    a = [tuple(pair.split(' ')) for pair in func_arg_str.split(', ')]
    try:
        ARG_TYPE = dict(a)
    except Exception as e:
        print(str(e))
        print(a)
    func_ret_str = get_return(func_name)
    a = [tuple(pair.split(' ')) for pair in func_ret_str.split(', ')]
    try:
        RET_TYPE = dict(a)
    except:
        RET_TYPE = {}

    if syntax.body.type == 'BlockStatement':
        for idx, block in enumerate(syntax.body.body):
            result_code.append(call_func_by_syntax(block, {'return': func_ret_str}))
            if idx == FUNC_LINES:
                break

    str_result_code = '\n'.join(result_code)

    panic_deal = ''
    if func_name.lower() in PANIC_DEAL_FUNC:
        panic_deal = '''defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
        '''

    return f'''func (self *{EX_NAME}) {capital_first(func_name)} ({func_arg_str}) ({func_ret_str}) {{
    {panic_deal}{str_result_code}
}}
    '''


def syntax_analysis(syntax, info={}):
    ret = ''
    for idx, func in enumerate(syntax):
        ret += call_func_by_syntax(func)
    return ret

def parse_by_syntax(str_code):
    str_code = str_code.replace('await ', '')
    str_code = str_code.replace('async ', 'function ')
    syntax = esprima.parse(str_code)
    return syntax_analysis(syntax.body)


def format_header():
    return f'''
    package {EX_NAME.lower()}

import (
	"fmt"
	. "github.com/georgexdz/ccxt/go/base"
	"reflect"
	"strings"
)

type {EX_NAME} struct {{
	Exchange
}}

func New(config *ExchangeConfig) (ex *{EX_NAME}, err error) {{
	ex = new({EX_NAME})
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {{
		ex = nil
		return
	}}

	return
}}

'''


def format_funcs(func_info_map):
    ret = ''

    for func, code in func_info_map.items():
        try:
            ret += parse_by_syntax(code) + '\n'
            # print(format_describe_func(info['describe']))
        except Exception as e:
            print(code)
            print(traceback.format_exc())

    return ret


def format_ex_code(ex):
    global CODE_INFO
    info = read_code_str(ex)
    CODE_INFO = info

    return f'''
    {format_header()}
    {format_describe_func(info['describe'])}
    {format_funcs(info['func'])}
    '''


def format_test_file(ex):
    return f'''
package {EX_NAME.lower()}

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func get_test_config(ex *{EX_NAME.capitalize()}) {{
	plan, err := ioutil.ReadFile("test_config.json")
	if err != nil {{
		return
	}}

	var data interface{{}}
	err = json.Unmarshal(plan, &data)
	if err != nil {{
		return
	}}
	
	fmt.Println(data)

	if json_config, ok := data.(map[string]interface{{}}); ok {{
        ex.Urls = map[string]interface{{}}{{
        	"api": map[string]interface{{}}{{
        		"public": json_config["url"],
        		"private": json_config["url"],
			}},
        }}
		ex.ApiUrls["private"] = json_config["url"].(string)
		ex.ApiUrls["public"] = json_config["url"].(string)
		ex.ApiKey = json_config["key"].(string)
		ex.Secret = json_config["secret"].(string)
		ex.Password = json_config["password"].(string)
	}}
}}

func TestFetchOrderBook(t *testing.T) {{
	ex, _ := New(nil)
	fmt.Println(ex.ApiDecodeInfo)
	ex.Verbose = true

	get_test_config(ex)

	markets, err := ex.LoadMarkets()
	if err != nil {{
		t.Fatal(err)
		return
	}}
	fmt.Println("markets:", markets)

	orderbook, err := ex.FetchOrderBook("BTC/USDT", 20, nil)
	if err != nil {{
		fmt.Println(err.Error())
		return
	}}
	fmt.Println("orderbook:", orderbook)

	ex.FetchBalance(nil)

	order, err := ex.CreateOrder("ETH/BTC", "limit", "buy", 0.0001, 0.024, nil)
	if err != nil {{
		return
	}}

	fmt.Println(ex.FetchOrder(order["id"].(string), "ETH/BTC", nil))

	openOrders, err := ex.FetchOpenOrders("ETH/BTC", 0, 20, nil)
	if err == nil {{
		fmt.Println("openorders", openOrders)
	}}

	if err == nil {{
		res, err := ex.CancelOrder(order["id"].(string), "ETH/BTC", nil)
		fmt.Println(res, err)
	}}
}}

//func main() {{
//	ex := &ccxt.Kucoin{{}}
//	ex.Init()
//	// testFetchMarkets(ex)
//	fmt.Println("enter")
//	testFetchOrderBook(ex)
//}}'''

def write_ex_file(ex, code):
    global PWD
    des_dir = os.path.join(PWD, '..', 'go', 'generated', f'{ex.lower()}')
    if not os.path.exists(des_dir):
        os.makedirs(des_dir)
    with open(os.path.join(des_dir, f'{ex.lower()}.go'), 'w') as f:
        f.write(code)
    with open(os.path.join(des_dir, f'{ex.lower()}_test.go'), 'w') as f:
        f.write(format_test_file(ex))
    # go fmt
    # cmd = "go fmt -x %s" % shlex.quote(des_dir)
    cmd = 'GO111MODULE=on go fmt -x %s' % shlex.quote(des_dir)
    p = subprocess.Popen(cmd, shell=True)
    p.communicate()
    if p.poll() == 0:
        print(ex, 'translate done:', des_dir)


def main():
    global EX_NAME
    for ex in sys.argv[1:]:
        EX_NAME = ex.capitalize()
        code = format_ex_code(ex)
        write_ex_file(ex, code)


if __name__ == '__main__':
    sys.exit(main())
