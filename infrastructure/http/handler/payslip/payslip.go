package payslip

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/entities"
	"github.com/riskykurniawan15/payrolls/infrastructure/http/middleware"
	"github.com/riskykurniawan15/payrolls/models/payslip"
	payslipService "github.com/riskykurniawan15/payrolls/services/payslip"
	"github.com/riskykurniawan15/payrolls/utils/logger"
)

type (
	IPayslipHandler interface {
		List(ctx echo.Context) error
		Generate(ctx echo.Context) error
		Print(ctx echo.Context) error
		GenerateSummary(ctx echo.Context) error
		PrintSummary(ctx echo.Context) error
	}

	PayslipHandler struct {
		logger          logger.Logger
		payslipServices payslipService.IPayslipService
	}
)

func NewPayslipHandlers(
	logger logger.Logger,
	payslipServices payslipService.IPayslipService,
) IPayslipHandler {
	return &PayslipHandler{
		logger:          logger,
		payslipServices: payslipServices,
	}
}

func (handler PayslipHandler) List(ctx echo.Context) error {
	requestID := middleware.GetRequestID(ctx)
	userID := middleware.GetUserID(ctx)

	// Parse query parameters
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	sortBy := ctx.QueryParam("sort_by")
	sortDesc := ctx.QueryParam("sort_desc") == "true"

	handler.logger.InfoT("incoming payslip list request", requestID, map[string]interface{}{
		"user_id":   userID,
		"page":      page,
		"limit":     limit,
		"sort_by":   sortBy,
		"sort_desc": sortDesc,
	})

	// Build request
	req := payslip.PayslipListRequest{
		Page:     page,
		Limit:    limit,
		SortBy:   sortBy,
		SortDesc: sortDesc,
	}

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.payslipServices.List(serviceCtx, req, userID)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return ctx.JSON(http.StatusInternalServerError, entities.ResponseFormater(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	handler.logger.InfoT("payslip list retrieved successfully", requestID, map[string]interface{}{
		"user_id": userID,
		"total":   response.Pagination.Total,
	})

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response.Data,
		"meta": response.Pagination,
	}))
}

func (handler PayslipHandler) Generate(ctx echo.Context) error {
	requestID := middleware.GetRequestID(ctx)
	userID := middleware.GetUserID(ctx)

	// Get period detail ID from URL parameter
	periodDetailIDStr := ctx.Param("id")
	periodDetailID, err := strconv.ParseUint(periodDetailIDStr, 10, 32)
	if err != nil {
		handler.logger.ErrorT("invalid period detail ID", requestID, map[string]interface{}{
			"user_id":          userID,
			"period_detail_id": periodDetailIDStr,
			"error":            "Invalid period detail ID format",
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid period detail ID",
		}))
	}

	// Get host from request
	host := ctx.Request().Host
	if ctx.Request().TLS != nil {
		host = "https://" + host
	} else {
		host = "http://" + host
	}

	handler.logger.InfoT("incoming generate payslip request", requestID, map[string]interface{}{
		"user_id":          userID,
		"period_detail_id": periodDetailID,
		"host":             host,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.payslipServices.GeneratePayslip(serviceCtx, uint(periodDetailID), userID, host)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"user_id":          userID,
			"period_detail_id": periodDetailID,
			"error":            err.Error(),
		})
		return ctx.JSON(http.StatusInternalServerError, entities.ResponseFormater(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	handler.logger.InfoT("payslip generated successfully", requestID, map[string]interface{}{
		"user_id":          userID,
		"period_detail_id": periodDetailID,
		"expires_at":       response.ExpiresAt,
	})

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}

func (handler PayslipHandler) Print(ctx echo.Context) error {
	requestID := middleware.GetRequestID(ctx)

	// Get token from query parameter
	token := ctx.QueryParam("token")
	if token == "" {
		handler.logger.ErrorT("missing token", requestID, map[string]interface{}{
			"error": "Token is required",
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Token is required",
		}))
	}

	handler.logger.InfoT("incoming print payslip request", requestID, map[string]interface{}{
		"token": token,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Get payslip data
	payslipData, err := handler.payslipServices.GetPayslipData(serviceCtx, token)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"token": token,
			"error": err.Error(),
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	// Generate HTML
	html, err := handler.generatePayslipHTML(payslipData)
	if err != nil {
		handler.logger.ErrorT("failed to generate HTML", requestID, map[string]interface{}{
			"token": token,
			"error": err.Error(),
		})
		return ctx.JSON(http.StatusInternalServerError, entities.ResponseFormater(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to generate payslip HTML",
		}))
	}

	handler.logger.InfoT("payslip HTML generated successfully", requestID, map[string]interface{}{
		"token": token,
	})

	// Return HTML response
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	return ctx.HTML(http.StatusOK, html)
}

func (handler PayslipHandler) GenerateSummary(ctx echo.Context) error {
	requestID := middleware.GetRequestID(ctx)

	// Get period ID from URL parameter
	periodIDStr := ctx.Param("id")
	periodID, err := strconv.ParseUint(periodIDStr, 10, 32)
	if err != nil {
		handler.logger.ErrorT("invalid period ID", requestID, map[string]interface{}{
			"period_id": periodIDStr,
			"error":     err.Error(),
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid period ID",
		}))
	}

	handler.logger.InfoT("incoming payslip summary generate request", requestID, map[string]interface{}{
		"period_id": periodID,
	})

	// Get host from request
	host := ctx.Scheme() + "://" + ctx.Request().Host

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	response, err := handler.payslipServices.GeneratePayslipSummary(serviceCtx, uint(periodID), host)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"period_id": periodID,
			"error":     err.Error(),
		})
		return ctx.JSON(http.StatusInternalServerError, entities.ResponseFormater(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	handler.logger.InfoT("payslip summary generated successfully", requestID, map[string]interface{}{
		"period_id": periodID,
		"print_url": response.PrintURL,
	})

	return ctx.JSON(http.StatusOK, entities.ResponseFormater(http.StatusOK, map[string]interface{}{
		"data": response,
	}))
}

func (handler PayslipHandler) PrintSummary(ctx echo.Context) error {
	requestID := middleware.GetRequestID(ctx)

	// Get token from query parameter
	token := ctx.QueryParam("token")
	if token == "" {
		handler.logger.ErrorT("missing token", requestID, map[string]interface{}{
			"error": "Token is required",
		})
		return ctx.JSON(http.StatusBadRequest, entities.ResponseFormater(http.StatusBadRequest, map[string]interface{}{
			"error": "Token is required",
		}))
	}

	handler.logger.InfoT("incoming payslip summary print request", requestID, map[string]interface{}{
		"token": token,
	})

	// Add request ID to context
	serviceCtx := middleware.AddRequestIDToContext(ctx.Request().Context(), requestID)

	// Call service
	data, err := handler.payslipServices.GetPayslipSummaryData(serviceCtx, token)
	if err != nil {
		handler.logger.ErrorT("service error", requestID, map[string]interface{}{
			"error": err.Error(),
		})
		return ctx.JSON(http.StatusInternalServerError, entities.ResponseFormater(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	// Generate HTML
	html, err := handler.generatePayslipSummaryHTML(data)
	if err != nil {
		handler.logger.ErrorT("failed to generate HTML", requestID, map[string]interface{}{
			"error": err.Error(),
		})
		return ctx.JSON(http.StatusInternalServerError, entities.ResponseFormater(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to generate payslip",
		}))
	}

	handler.logger.InfoT("payslip summary printed successfully", requestID, map[string]interface{}{
		"period_name": data.PeriodName,
	})

	return ctx.HTML(http.StatusOK, html)
}

func (handler PayslipHandler) generatePayslipSummaryHTML(data *payslip.PayslipSummaryData) (string, error) {
	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Management Payslip Summary</title>
  <style>
    @media print {
      body {
        margin: 0;
        -webkit-print-color-adjust: exact;
      }
    }

    body {
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      max-width: 800px;
      margin: 40px auto;
      padding: 20px;
    }

    h1, h2 {
      text-align: center;
      margin-bottom: 10px;
    }

    .section-title {
      background-color: #f0f0f0;
      padding: 8px;
      font-weight: bold;
      margin-top: 20px;
    }

    table {
      width: 100%;
      border-collapse: collapse;
      margin-top: 10px;
    }

    th, td {
      padding: 8px;
      border: 1px solid #ccc;
      text-align: left;
    }

    .right {
      text-align: right;
    }

    .total-row {
      font-weight: bold;
      background-color: #f9f9f9;
    }

    .footer {
      margin-top: 40px;
      text-align: center;
      font-size: 12px;
      color: #666;
    }
  </style>
</head>
<body onload="window.print()">

  <h1>{{.CompanyName}}</h1>
  <h2>Management Payslip Summary</h2>

  <p><strong>Period:</strong> {{.PeriodName}}</p>

  <div class="section-title">Summary</div>
  <table>
    <tr>
      <th>Description</th>
      <th class="right">Amount</th>
    </tr>
    <tr>
      <td>Total Employees</td>
      <td class="right">{{.TotalEmployees}}</td>
    </tr>
    <tr class="total-row">
      <td>Total Take Home Pay</td>
      <td class="right">{{formatRupiah .TotalTakeHomePay}}</td>
    </tr>
  </table>

  <div class="section-title">Employee Take Home Pay List</div>
  <table>
    <tr>
      <th>No</th>
      <th>Employee Name</th>
      <th class="right">THP</th>
    </tr>
    {{range .EmployeeList}}
    <tr>
      <td>{{.No}}</td>
      <td>{{.EmployeeName}}</td>
      <td class="right">{{formatRupiah .TakeHomePay}}</td>
    </tr>
    {{end}}
    <tr class="total-row">
      <td colspan="2">Total</td>
      <td class="right">{{formatRupiah .TotalTakeHomePay}}</td>
    </tr>
  </table>

  <div class="footer">
    This report was generated automatically and does not require a signature.
  </div>

</body>
</html>`

	// Create template with custom function
	funcMap := template.FuncMap{
		"formatRupiah": handler.formatRupiah,
	}

	tmpl, err := template.New("payslip_summary").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

func (handler PayslipHandler) formatRupiah(amount float64) string {
	// Convert to integer (in rupiah)
	rupiah := int64(amount)

	// Get decimal part (2 digits)
	decimal := int64((amount - float64(rupiah)) * 100)

	// Format with thousand separators
	formatted := fmt.Sprintf("%d", rupiah)

	// Add thousand separators
	var result strings.Builder
	length := len(formatted)
	for i, digit := range formatted {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteString(".")
		}
		result.WriteRune(digit)
	}

	// Add decimal part
	return "Rp. " + result.String() + "," + fmt.Sprintf("%02d", decimal)
}

func (handler PayslipHandler) generatePayslipHTML(data *payslip.PayslipData) (string, error) {
	// Create template functions
	funcMap := template.FuncMap{
		"formatRupiah": handler.formatRupiah,
	}

	const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Payslip</title>
  <style>
    @media print {
      body {
        margin: 0;
        -webkit-print-color-adjust: exact;
      }
    }

    body {
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      max-width: 800px;
      margin: 40px auto;
      padding: 20px;
    }

    h1, h2 {
      text-align: center;
      margin-bottom: 10px;
    }

    .section-title {
      background-color: #f0f0f0;
      padding: 8px;
      font-weight: bold;
      margin-top: 20px;
    }

    table {
      width: 100%;
      border-collapse: collapse;
      margin-top: 10px;
    }

    th, td {
      padding: 8px;
      border: 1px solid #ccc;
      text-align: left;
    }

    .right {
      text-align: right;
    }

    .total-row {
      font-weight: bold;
      background-color: #f9f9f9;
    }

    .footer {
      margin-top: 40px;
      text-align: center;
      font-size: 12px;
      color: #666;
    }
  </style>
</head>
<body onload="window.print()">

  <h1>{{.CompanyName}}</h1>
  <h2>Payslip</h2>

  <p><strong>Employee:</strong> {{.EmployeeName}}<br>
     <strong>Period:</strong> {{.PeriodName}}</p>

  <div class="section-title">Work Summary</div>
  <table>
    <tr>
      <th>Description</th>
      <th class="right">Amount</th>
    </tr>
    <tr>
      <td>Total Working Days</td>
      <td class="right">{{.TotalWorking}}</td>
    </tr>
    <tr>
      <td>Daily Rate</td>
      <td class="right">{{formatRupiah .DailyRate}}</td>
    </tr>
    <tr class="total-row">
      <td>Base Salary</td>
      <td class="right">{{formatRupiah .BaseSalary}}</td>
    </tr>
  </table>

  {{if .OvertimeDetails}}
  <div class="section-title">Overtime</div>
  <table>
    <tr>
      <th>Date</th>
      <th>Hours</th>
      <th class="right">Amount</th>
    </tr>
    {{range .OvertimeDetails}}
    <tr>
      <td>{{.Date}}</td>
      <td>{{.Hours}}</td>
      <td class="right">{{formatRupiah .Amount}}</td>
    </tr>
    {{end}}
    <tr class="total-row">
      <td colspan="2">Total Overtime</td>
      <td class="right">{{formatRupiah .TotalOvertime}}</td>
    </tr>
  </table>
  {{end}}
  
  {{if .Reimbursements}}
  <div class="section-title">Reimbursements</div>
  <table>
    <tr>
      <th>Description</th>
      <th class="right">Amount</th>
    </tr>
    {{range .Reimbursements}}
    <tr>
      <td>{{.Title}}</td>
      <td class="right">{{formatRupiah .Amount}}</td>
    </tr>
    {{end}}
    <tr class="total-row">
      <td>Total Reimbursement</td>
      <td class="right">{{formatRupiah .TotalReimbursement}}</td>
    </tr>
  </table>
  {{end}}

  <div class="section-title">Net Salary</div>
  <table>
    <tr class="total-row">
      <td>Take Home Pay</td>
      <td class="right">{{formatRupiah .TakeHomePay}}</td>
    </tr>
  </table>

  <div class="footer">
    This payslip was generated automatically and does not require a signature.<br>
    Generated at: {{.GeneratedAt.Format "2006-01-02 15:04:05"}}
  </div>

</body>
</html>`

	tmpl, err := template.New("payslip").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	// Execute template
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
