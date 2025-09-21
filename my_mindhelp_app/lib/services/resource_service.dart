import '../core/api_client.dart';
import '../core/api_config.dart';
import '../models/resource.dart';
import '../models/counselor.dart';

class ResourceService {
  final ApiClient _apiClient = ApiClient();

  /// 獲取諮商師列表
  Future<CounselorListResponse> getCounselors({
    String? search,
    String? workLocation,
    String? specialty,
    int page = 1,
    int pageSize = 10,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };

    if (search != null && search.isNotEmpty) {
      queryParams['search'] = search;
    }
    if (workLocation != null && workLocation.isNotEmpty) {
      queryParams['work_location'] = workLocation;
    }
    if (specialty != null && specialty.isNotEmpty) {
      queryParams['specialty'] = specialty;
    }

    final response = await _apiClient.dio.get(
      ApiConfig.counselors,
      queryParameters: queryParams,
    );

    if (response.statusCode == 200) {
      return CounselorListResponse.fromJson(response.data['data']);
    }
    throw Exception('獲取諮商師列表失敗');
  }

  /// 獲取單一諮商師詳情
  Future<Counselor> getCounselor(String counselorId) async {
    final response = await _apiClient.dio.get('${ApiConfig.counselors}/$counselorId');

    if (response.statusCode == 200) {
      return Counselor.fromJson(response.data['data']);
    }
    throw Exception('獲取諮商師詳情失敗');
  }

  /// 獲取諮商所列表
  Future<CounselingCenterListResponse> getCounselingCenters({
    String? search,
    String? location,
    int page = 1,
    int pageSize = 10,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };

    if (search != null && search.isNotEmpty) {
      queryParams['search'] = search;
    }
    if (location != null && location.isNotEmpty) {
      queryParams['location'] = location;
    }

    final response = await _apiClient.dio.get(
      ApiConfig.counselingCenters,
      queryParameters: queryParams,
    );

    if (response.statusCode == 200) {
      return CounselingCenterListResponse.fromJson(response.data['data']);
    }
    throw Exception('獲取諮商所列表失敗');
  }

  /// 獲取單一諮商所詳情
  Future<CounselingCenter> getCounselingCenter(String centerId) async {
    final response = await _apiClient.dio.get('${ApiConfig.counselingCenters}/$centerId');

    if (response.statusCode == 200) {
      return CounselingCenter.fromJson(response.data['data']);
    }
    throw Exception('獲取諮商所詳情失敗');
  }

  /// 獲取推薦醫師列表
  Future<RecommendedDoctorListResponse> getRecommendedDoctors({
    String? search,
    String? specialty,
    String? location,
    int page = 1,
    int pageSize = 10,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };

    if (search != null && search.isNotEmpty) {
      queryParams['search'] = search;
    }
    if (specialty != null && specialty.isNotEmpty) {
      queryParams['specialty'] = specialty;
    }
    if (location != null && location.isNotEmpty) {
      queryParams['location'] = location;
    }

    final response = await _apiClient.dio.get(
      ApiConfig.recommendedDoctors,
      queryParameters: queryParams,
    );

    if (response.statusCode == 200) {
      return RecommendedDoctorListResponse.fromJson(response.data['data']);
    }
    throw Exception('獲取推薦醫師列表失敗');
  }

  /// 獲取單一推薦醫師詳情
  Future<RecommendedDoctor> getRecommendedDoctor(String doctorId) async {
    final response = await _apiClient.dio.get('${ApiConfig.recommendedDoctors}/$doctorId');

    if (response.statusCode == 200) {
      return RecommendedDoctor.fromJson(response.data['data']);
    }
    throw Exception('獲取推薦醫師詳情失敗');
  }

  /// 獲取地圖地址資訊
  Future<Map<String, dynamic>> getMapAddresses({
    String? type,
    int? limit,
  }) async {
    final queryParams = <String, dynamic>{};

    if (type != null && type.isNotEmpty) {
      queryParams['type'] = type;
    }
    if (limit != null) {
      queryParams['limit'] = limit;
    }

    final response = await _apiClient.dio.get(
      '${ApiConfig.maps}/addresses',
      queryParameters: queryParams,
    );

    if (response.statusCode == 200) {
      return response.data['data'];
    }
    throw Exception('獲取地圖地址失敗');
  }

  /// 獲取 Google Maps 專用格式地址
  Future<Map<String, dynamic>> getGoogleMapAddresses({
    String? format,
  }) async {
    final queryParams = <String, dynamic>{};

    if (format != null && format.isNotEmpty) {
      queryParams['format'] = format;
    }

    final response = await _apiClient.dio.get(
      '${ApiConfig.maps}/google-addresses',
      queryParameters: queryParams,
    );

    if (response.statusCode == 200) {
      return response.data['data'];
    }
    throw Exception('獲取 Google Maps 地址失敗');
  }

  /// 收藏資源
  Future<void> bookmarkResource(String resourceId) async {
    final response = await _apiClient.dio.post(
      '${ApiConfig.resources}/$resourceId/bookmark',
    );

    if (response.statusCode != 204 && response.statusCode != 200) {
      throw Exception('收藏資源失敗');
    }
  }

  /// 取消收藏資源
  Future<void> unbookmarkResource(String resourceId) async {
    final response = await _apiClient.dio.delete(
      '${ApiConfig.resources}/$resourceId/bookmark',
    );

    if (response.statusCode != 204 && response.statusCode != 200) {
      throw Exception('取消收藏資源失敗');
    }
  }

  /// 獲取使用者收藏的資源列表
  Future<ResourceListResponse> getBookmarkedResources({
    int page = 1,
    int limit = 10,
  }) async {
    final response = await _apiClient.dio.get(
      '${ApiConfig.users}/me/bookmarks/resources',
      queryParameters: {
        'page': page,
        'limit': limit,
      },
    );

    if (response.statusCode == 200) {
      return ResourceListResponse.fromJson(response.data['data']);
    }
    throw Exception('獲取收藏資源列表失敗');
  }
}
